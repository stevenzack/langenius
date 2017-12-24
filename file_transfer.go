package LanGenius

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
)

type FileEntry struct {
	Name, Path string
}

var (
	homeData struct {
		Files            []FileEntry
		Clipboard        string
		ClipboardEnabled bool
	}
)

func home(w http.ResponseWriter, r *http.Request) {
	t := template.New("homeTPL")
	t.Parse(`
<!DOCTYPE html>
<html>
<head>
<title>LanGenius - Transfer file easily</title>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<script>function initMe(){if("zh-CN"==navigator.language){null!=document.getElementById("btCopy")&&(document.getElementById("btCopy").value="复制"),null!=document.getElementById("btSend")&&(document.getElementById("btSend").value="发送"),null!=document.getElementById("txtClipboard")&&(document.getElementById("txtClipboard").innerHTML="剪切板"),null!=document.getElementById("txtBindClip")&&(document.getElementById("txtBindClip").innerHTML="同步"),document.getElementById("txtFiles").innerHTML="文件",document.getElementById("liveStreamPage").innerHTML="直播页面",window.Succeed="成功",document.getElementById("submit_button").value="上传文件";for(var e=document.getElementsByClassName("downloadA"),t=0;t<e.length;t++)e[t].innerHTML="下载"}else window.Succeed=" Succeed";"checked"==getCookie("LanGenius-Bind-Clipboard")&&(document.getElementById("btBind").checked=!0,BindClipboard()),detectLink()}function detectLink(){if(null!=document.getElementById("cb")){var e=document.getElementById("cb").value,t="";if(e.indexOf("http://")>-1){var n=e.indexOf(" ",e.indexOf("http://"));t=e.substring(e.indexOf("http://"),-1==n?e.length:n)}else if(e.indexOf("https://")>-1){n=e.indexOf(" ",e.indexOf("https://"));t=e.substring(e.indexOf("https://"),-1==n?e.length:n)}""!=t?("zh-CN"==navigator.language?document.getElementById("openUrl").innerHTML="打开网址":document.getElementById("openUrl").innerHTML="Open URL",document.getElementById("openUrl").href=t):document.getElementById("openUrl").innerHTML=""}}function setCookie(e,t,n){var d=new Date;d.setTime(d.getTime()+24*n*60*60*1e3);var o="expires="+d.toGMTString();document.cookie=e+"="+t+"; "+o}function getCookie(e){for(var t=e+"=",n=document.cookie.split(";"),d=0;d<n.length;d++){var o=n[d].trim();if(0==o.indexOf(t))return o.substring(t.length,o.length)}return""}function copy(){detectLink();var e=document.getElementById("btCopy");document.getElementById("cb").select(),document.execCommand("Copy");var t=e.value;e.value=window.Succeed,e.disabled="disabled",setTimeout("document.getElementById('btCopy').value='"+t+"';document.getElementById('btCopy').disabled=''",1e3)}function send(){var e=document.getElementById("btSend"),t=e.value,n=document.getElementById("cb").value,d=new XMLHttpRequest,o=new FormData,l=new Object;l.Cb=n,o.append("cb",JSON.stringify(l)),d.onload=function(n){200!=this.status&&302!=this.status&&304!=this.status||(e.value=window.Succeed,e.disabled="disabled",setTimeout("document.getElementById('btSend').value='"+t+"';document.getElementById('btSend').disabled='';",1e3))},d.open("POST","/send",!0),d.send(o),console.log("sent clip")}function DoUpload(){for(var e=new FormData,t=document.getElementById("myUploadFile"),n=0;n<t.files.length;n++)e.append("myUploadFile",t.files[n]);uploading();var d=new XMLHttpRequest;d.upload.addEventListener("progress",function(e){if(e.lengthComputable){var t=Math.round(100*e.loaded/e.total);document.getElementById("uploadProgress").value=t,document.getElementById("result").innerHTML=t.toString()+"%"}else document.getElementById("result").innerHTML="unable to compute"},!1),d.onreadystatechange=function(){4==d.readyState&&200==d.status&&(document.getElementById("result").innerHTML=d.responseText,setTimeout("uploadDone()",1e3))},d.addEventListener("load",function(e){document.getElementById("result").innerHTML=e.target.responseText,setTimeout("uploadDone()",1e3)},!1),d.addEventListener("error",function(e){document.getElementById("result").innerHTML="Something went wrong .",setTimeout("uploadDone()",1e3)},!1),d.addEventListener("abort",function(e){document.getElementById("result").innerHTML="Abort .",setTimeout("uploadDone()",1e3)},!1),d.open("POST","/upload"),d.send(e)}function uploading(){"zh-CN"==navigator.language?document.getElementById("submit_button").value="上传中...":document.getElementById("submit_button").value="Uploading",document.getElementById("submit_button").disabled="disabled",document.getElementById("myUploadFile").disabled="disabled"}function uploadDone(){document.getElementById("uploadProgress").value=0,"zh-CN"==navigator.language?document.getElementById("submit_button").value="上传文件":document.getElementById("submit_button").value="Upload",document.getElementById("submit_button").disabled="",document.getElementById("myUploadFile").disabled="",location.href=location.href}function BindClipboard(){document.getElementById("btBind").checked?(setCookie("LanGenius-Bind-Clipboard","checked",365),window.wsBind=new WebSocket("ws://"+location.href.split("/")[2]+"/wsClipboard"),window.wsBind.onmessage=function(e){var t=JSON.parse(e.data);document.getElementById("cb").value=t.Content,copy()},window.wsBind.onopen=function(){console.log("bind open")},window.wsBind.onerror=function(e){console.log("bind err:"+e.data)},window.wsBind.onclose=function(){console.log("bind closed"),window.wsBind=null}):(setCookie("LanGenius-Bind-Clipboard","unchecked",365),null!=window.wsBind&&(window.wsBind.close(),window.wsBind=null))}setTimeout("initMe()",50)</script>
<style>.wrapper{background-color:#fff;display:inline-block;padding:5px;box-shadow:2px 2px 10px #000;border-radius:10px}</style>
</head>
<body style="background-color:#58c6d5;font-family:sans-serif;padding:10px;margin:0">
<center>
{{if .ClipboardEnabled}}
<div class="wrapper">
<table>
<tr>
<th style="color:#d81b60" id="txtClipboard">Clipboard</th>
</tr>
<tr align="center">
<td>
<textarea name="cb" id="cb" cols="30" rows="5" onkeyup="detectLink()">{{.Clipboard}}</textarea>
</td>
<td>
<a href="" id="openUrl" target="_blank"></a>
<br>
<input type="button" value="Copy" id="btCopy" onclick="copy(this)">
<br>
<input type="button" value="Send" id="btSend" onclick="send(this)">
<br>
<input type="checkbox" id="btBind" onchange="BindClipboard()"><span id="txtBindClip">Bind</span>
</td>
<td><span id="spanInfo"></span>
<br><span></span></td>
</tr>
<tr>
<td colspan="2">
<hr>
</td>
</tr>
</table>
</div>
{{end}}
<br>
<br>
<div class="wrapper">
<table>
<tr>
<th style="color:#1e88e5" id="txtFiles">Files</th>
</tr>
<tr>
<td colspan="2">
{{range .Files}}
<a href="/viewfile/{{.Name}}">{{.Name}}</a>
<a href="/download/{{.Name}}" class="downloadA">Download</a>
<br>{{end}}
</td>
</tr>
<tr>
<td colspan="2">
<hr>
</td>
</tr>
<tr>
<td>
<input type="file" name="myUploadFile" id="myUploadFile" onchange='document.getElementById("submit_button").disabled=""' multiple>
<input type="button" id="submit_button" disabled value="Upload" onclick="DoUpload()">
</td>
</tr>
<tr>
<td>
<progress value="0" max="100" id="uploadProgress" style="height:4px;width:100%"></progress>
<div id="result"></div>
</td>
</tr>
<tr>
<td>
<a href="/live" id="liveStreamPage">Live Page</a>
</td>
</tr>
</table>
</div>
<br>
</center>
</body>
</html>
		`)
	t.Execute(w, homeData)
	// t, e := template.ParseFiles("/home/asd/go/src/LanGenius/views/index.html")
	// if e != nil {
	// 	fmt.Println(e)
	// 	return
	// }
	// t.Execute(w, homeData)
}
func send(w http.ResponseWriter, r *http.Request) {
	var gobj struct {
		Cb string
	}
	e := json.Unmarshal([]byte(r.FormValue("cb")), &gobj)
	if e != nil {
		fmt.Fprint(w, e.Error())
		return
	}
	mEventHandler.OnClipboardReceived(gobj.Cb, true)
	homeData.Clipboard = gobj.Cb
	fmt.Fprint(w, "OK")
}
func download(w http.ResponseWriter, r *http.Request) {
	filename, _ := url.QueryUnescape(r.URL.RequestURI()[len("/download/"):])
	for _, v := range homeData.Files {
		if v.Name == filename {
			if filename[len(filename)-4:] == ".apk" {
				w.Header().Set("Content-Type", "application/vnd.android.package-archive")
			}
			w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
			http.ServeFile(w, r, v.Path)
			return
		}
	}
	http.NotFound(w, r)
}
func viewfile(w http.ResponseWriter, r *http.Request) {
	filename, _ := url.QueryUnescape(r.URL.RequestURI()[len("/viewfile/"):])
	for _, v := range homeData.Files {
		if v.Name == filename {
			if filename[len(filename)-4:] == ".apk" {
				w.Header().Set("Content-Type", "application/vnd.android.package-archive")
			}
			http.ServeFile(w, r, v.Path)
			return
		}
	}
	http.NotFound(w, r)
}
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	fhs := r.MultipartForm.File["myUploadFile"]
	for _, v := range fhs {
		file, e := v.Open()
		if e != nil {
			fmt.Println(e)
			fmt.Fprint(w, e.Error())
			return
		}
		mEventHandler.OnFileReceived(v.Filename)
		mf, e := os.OpenFile(storagePath+v.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if e != nil {
			fmt.Println(e)
			fmt.Fprint(w, e.Error())
			return
		}
		defer mf.Close()
		io.Copy(mf, file)
	}
	fmt.Fprint(w, "OK")
}
func AddFile(str string) {
	if !contains(str, "/") {
		return
	}
	homeData.Files = append(homeData.Files, FileEntry{Name: getFileName(str), Path: str})
}
func RemoveFile(index int) {
	if index > -1 && index < len(homeData.Files) {
		homeData.Files = append(homeData.Files[:index], homeData.Files[index+1:]...)
	} else {
		fmt.Println("remove file index out of bound: index=", index)
	}
}
func SetStoragePath(str string) {
	storagePath = str
}

//clipboard part

var (
	clipConns []*websocket.Conn
)

func SetClipboardEnabled(b bool) {
	homeData.ClipboardEnabled = b
}

func SetClipboard(str string) {
	homeData.Clipboard = str

	if len(clipConns) > 0 {
		b, _ := json.Marshal(Msg{Content: str, State: "OK"})
		for k, v := range clipConns {
			e := websocket.Message.Send(v, string(b))
			if e != nil {
				fmt.Println("ws send failed:", e.Error())
				clipConns = append(clipConns[:k], clipConns[k+1:]...)
			}
		}
	}
}

func wsClipboard(ws *websocket.Conn) {
	defer ws.Close()
	clipConns = append(clipConns, ws)
	for {
		s := ""
		e := websocket.Message.Receive(ws, &s)
		if e != nil {
			fmt.Println("ws receive failed:", e.Error())
			for k, v := range clipConns {
				if v == ws {
					clipConns = append(clipConns[:k], clipConns[k+1:]...)
					return
				}
			}
			return
		}
		cm := Msg{}
		e = json.Unmarshal([]byte(s), &cm)
		if e != nil {
			fmt.Println("unmarshal failed :", e.Error(), s)
			return
		}
		mEventHandler.OnClipboardReceived(cm.Content, false)
		homeData.Clipboard = cm.Content
	}
}
func Logd(e interface{}) {
	fmt.Println(e)
}
