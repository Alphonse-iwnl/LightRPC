package protocol

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	. "rpc/LightProtoPlug/common"
)

type HttpTyper struct {
	ServiceName string
	MethodInfo  []MethodInfo
	Content     map[string]string
	FileName    string
}

// NewHttpTyper init
func NewHttpTyper(fileName string, serviceName string, Methods []MethodInfo) *HttpTyper {
	NewType := &HttpTyper{}
	NewType.Content = make(map[string]string)
	NewType.MethodInfo = Methods
	NewType.ServiceName = serviceName
	NewType.FileName = fileName
	return NewType
}

func (ht *HttpTyper) OutputServiceCode(mode string) bool {
	ht.ImportBuild()
	// common part
	PrintContent := ht.Content["IS"] + "\n"
	if mode == "cli" {
		PrintContent = ht.generateClientContent(PrintContent)
	} else if mode == "server" {
		PrintContent = ht.generateServerContent(PrintContent)
	} else if mode == "all" {
		PrintContent = ht.generateClientContent(PrintContent)
		PrintContent = ht.generateServerContent(PrintContent)
	} else {
		fmt.Printf("unknown mode:%s, cant generate code", mode)
		return false
	}
	err := _writeToFile(PrintContent, ht.FileName)
	if err != nil {
		fmt.Printf("write code to file:%s error:%v.", ht.FileName, err)
		return false
	}
	fmt.Println("write to file success.")
	return true
}

func (ht *HttpTyper) generateClientContent(content string) string {
	ht.ClientInterfaceBuild()
	ht.ClientRegisterBuild()
	ht.ClientHandlerBuild()
	content += ht.Content["CI"] + "\n" + ht.Content["CR"] + "\n" + ht.Content["CH"]
	return content
}

func (ht *HttpTyper) generateServerContent(content string) string {
	ht.ServiceDescBuild()
	ht.ServerInterfaceBuild()
	ht.ServerRegisterBuild()
	ht.ServerHandlerBuild()
	content += ht.Content["Desc"] + ht.Content["SI"] + "\n" + ht.Content["SR"] + "\n" + ht.Content["SH"]
	return content
}

// import build
func (ht *HttpTyper) ImportBuild() {
	ImportStr := "\nimport . \"net/http\"\n"
	ImportStr += "import \"io/ioutil\"\n"
	ImportStr += "import \"strings\"\n"
	ImportStr += "import \"encoding/json\"\n"
	ImportStr += "import \"context\"\n"

	ImportStr += "import model \"rpc/LightRpc/common\"\n"
	ImportStr += "\n"
	// ht.Content["IS"] = ImportStr

	protoFile, err := os.Open(ht.FileName)
	if err != nil {
		fmt.Printf("open file:%s error:%v", ht.FileName, err)
		panic("open file error.")
	}
	defer func() {
		_ = protoFile.Close()
	}()
	fileData, err := ioutil.ReadAll(protoFile)
	if err != nil {
		fmt.Printf("load file:%s content error:%v", ht.FileName, err)
		panic("read file error.")
	}

	fileStr := string(fileData)

	insertIndex := 0
	for index, r := range fileStr{
		if string(r) == "i" && fileStr[index:index+6] == "import"{
			insertIndex = index
			break
		}
	}
	newContent:=fileStr[:insertIndex]
	newContent +=ImportStr
	newContent +=fileStr[insertIndex:]

	ht.Content["IS"] = newContent
}

// ServiceDescBuild common part ServiceDesc
func (ht *HttpTyper) ServiceDescBuild() {
	DescStr := ""
	DescStr += "func (s *ServiceMiddle) ServiceDescBuild() model.ServiceDesc{\n"
	DescStr += "	var " + ht.ServiceName + "_ServiceDesc = model.ServiceDesc{\n"
	DescStr += "		ServiceName: \"" + ht.ServiceName + "\",\n"
	DescStr += "		HandlerType: (*" + ht.ServiceName + "Server)(nil),\n"
	DescStr += "		Methods: map[string]model.HttpHandler{\n"
	for _, item := range ht.MethodInfo {
		DescStr += "			\"" + item.MethodName + "\": s._" + item.MethodName + "_Handler,\n"
	}
	DescStr += "		},\n"
	DescStr += "		MetaData:\"" + ht.FileName + "\",\n"
	DescStr += "	}\n"
	DescStr += "	return " + ht.ServiceName + "_ServiceDesc\n"
	DescStr += "}\n"
	DescStr += "\n"
	ht.Content["Desc"] = DescStr
}

// ServerInterfaceBuild Server part
func (ht *HttpTyper) ServerInterfaceBuild() {
	SIStr := "type " + ht.ServiceName + "Server interface {\n"
	for _, method := range ht.MethodInfo {
		SIStr += "		" + method.MethodName + "(ctx context.Context, in *" + method.InputName + ") (*" + method.OutputName + ",error)\n"
	}
	SIStr += "}\n"
	ht.Content["SI"] = SIStr
}

func (ht *HttpTyper) ServerRegisterBuild() {
	SRStr := "func RegisterServiceServer (s ServiceMiddle, handler " + ht.ServiceName + "Server){\n"
	SRStr += "		s.S.RegisterService(s.ServiceDescBuild(), handler)\n"
	SRStr += "}\n"

	ht.Content["SR"] = SRStr
}

func (ht *HttpTyper) ServerHandlerBuild() {
	// ...http入口,调用method.methodinterface
	SHStr := "type ServiceMiddle struct {\n"
	SHStr += "	S *model.Server\n"
	SHStr += "}\n"

	for _, item := range ht.MethodInfo {
		SHStr += "func (s *ServiceMiddle) _" + item.MethodName + "_Handler(w ResponseWriter, r *Request){\n"

		// decode http.body(json or proto data) to message
		SHStr += "		in := new(" + item.InputName + ")\n"
		SHStr += "		body, err := ioutil.ReadAll(r.Body)\n"
		SHStr += "		if err != nil {\n"

		// write 400 to response
		SHStr += "			s.S.Errors(w, r, err, 400)\n"
		SHStr += "			return\n"
		SHStr += "		}\n"
		SHStr += "		form := r.Header.Get(\"Content-Type\")\n"
		SHStr += "		if strings.Contains(form,\"proto\"){\n"
		SHStr += "			err = proto.Unmarshal(body, in)\n"
		SHStr += "				if err != nil {\n"
		SHStr += "					s.S.Errors(w, r, err, 400)\n"
		SHStr += "					return\n"
		SHStr += "				}\n"
		SHStr += "		}else if strings.Contains(form, \"json\"){\n"
		SHStr += "			err = json.Unmarshal(body, in)\n"
		SHStr += "				if err != nil {\n"
		SHStr += "					s.S.Errors(w, r, err, 400)\n"
		SHStr += "					return\n"
		SHStr += "				}\n"
		SHStr += "		}\n"

		// run interface impl in Server struct
		SHStr += "		handler,ok := s.S.M[s.S.SD.ServiceName]\n"
		SHStr += "		if !ok{\n"
		SHStr += "			s.S.Errors(w, r, err, 500)\n"
		SHStr += "			return\n"
		SHStr += "		}\n"
		SHStr += "		response, err := handler.("+ht.ServiceName+"Server)."+item.MethodName+"(context.Background(), in)\n"
		SHStr += "		if err != nil {\n"
		SHStr += "			s.S.Errors(w, r, err, 500)\n"
		SHStr += "			return\n"
		SHStr += "		}\n"
		//SHStr += "		response := s.S.(" + ht.ServiceName + "Server)." + item.MethodName + "(context.Background(), in)\n"
		// code message(json or proto) to byte data
		SHStr += "		var bData []byte\n"
		SHStr += "		if strings.Contains(form,\"proto\"){\n"
		SHStr += "			bData, err = proto.Marshal(response)\n"
		SHStr += "				if err != nil {\n"
		SHStr += "					s.S.Errors(w, r, err, 500)\n"
		SHStr += "					return\n"
		SHStr += "				}\n"
		SHStr += "		}else if strings.Contains(form, \"json\"){\n"
		SHStr += "			bData, err = json.Marshal(response)\n"
		SHStr += "				if err != nil {\n"
		SHStr += "					s.S.Errors(w, r, err, 500)\n"
		SHStr += "					return\n"
		SHStr += "				}\n"
		SHStr += "		}\n"

		// write to response body
		SHStr += "		_ , err = w.Write(bData)\n"
		SHStr += "		if err != nil {\n"
		SHStr += "			s.S.Errors(w, r, err, 500)\n"
		SHStr += "			return\n"
		SHStr += "		}\n"
		SHStr += "		return\n"
		SHStr += "}\n"

		SHStr += "\n"
	}

	ht.Content["SH"] = SHStr
}


// ClientInterfaceBuild Client part
func (ht *HttpTyper) ClientInterfaceBuild() {
	CIStr := "type " + ht.ServiceName + "Client interface {\n"
	for _, method := range ht.MethodInfo {
		CIStr += "		" + method.MethodName + "(ctx context.Context, in *" + method.InputName + ") (*" + method.OutputName + ",error)\n"
	}
	CIStr += "}\n"
	ht.Content["CI"] = CIStr
}

func (ht *HttpTyper) ClientRegisterBuild() {
	CRStr := "type " + ht.ServiceName + "Struct struct {\n"
	CRStr += "		cs *model.ClientConn\n"
	CRStr += "}\n"

	CRStr += "\n\n"

	CRStr += "func New" + ht.ServiceName + "Client(cs *model.ClientConn) *" + ht.ServiceName + "Struct {\n"
	CRStr += "		return &" + ht.ServiceName + "Struct{cs}\n"
	CRStr += "}\n"

	ht.Content["CR"] = CRStr
}

// ClientHandlerBuild client start rpc request
func (ht *HttpTyper) ClientHandlerBuild() {
	CHStr := ""
	for _, item := range ht.MethodInfo {
		CHStr += "func (c *" + ht.ServiceName + "Struct) " + item.MethodName + " (ctx context.Context, in *" + item.InputName + ") (*" + item.OutputName + ", error){\n"
		CHStr += "		out := new(" + item.OutputName + ")\n"
		CHStr += "		inputBody ,err := proto.Marshal(in)\n"
		CHStr += "		if err != nil {\n"
		CHStr += "			return nil,err\n"
		CHStr += "		}\n"
		CHStr += "		outPutBody, err := c.cs.Invoke(ctx, \"" + ht.ServiceName + "/" + item.MethodName + "\", inputBody)\n"
		CHStr += "		if err != nil {\n"
		CHStr += "			return nil, err \n"
		CHStr += "		}\n"
		CHStr += "		err = proto.Unmarshal(outPutBody, out)\n"
		CHStr += "		if err != nil{\n"
		CHStr += "			return nil,err\n"
		CHStr += "		}\n"
		CHStr += "		return out, nil \n"
		CHStr += "}\n"
		CHStr += "\n"
	}

	ht.Content["CH"] = CHStr
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func _writeToFile(content string, fileName string) error {
	if checkFileIsExist(fileName) {
		//var f *os.File
		//var err error
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, 0666)
		defer func() {
			_ = f.Close()
		}()
		if err != nil {
			fmt.Println("Open .pb file fail.")
			return err
		}
		//_, err = f.Write([]byte(content))
		//w := bufio.NewWriter(f) //创建新的 Writer 对象
		//_, err = w.WriteString(content)
		_, err = io.WriteString(f,content)
		//_ =w.Flush()
		if err != nil {
			fmt.Println("Write content to .pb.go file error.")
			return err
		}
		//defer func() {
		//	_ = w.Flush()
		//}()
		return nil
	} else {
		fmt.Println(".pb file not exist, please check and try again.")
		return errors.New("file not found error")
	}
}
