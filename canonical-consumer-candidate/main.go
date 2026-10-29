package main
import (
 "encoding/json"
 "fmt"
 "os"
 "runtime"
 "runtime/debug"
 "strings"
 "github.com/PuerkitoBio/goquery"
)
func main() {
 input, err := os.ReadFile("input.html"); if err != nil { panic(err) }
 info, ok := debug.ReadBuildInfo(); if !ok { panic("missing build info") }
 module := ""; for _, d := range info.Deps { if d.Path == "github.com/PuerkitoBio/goquery" { if d.Replace != nil { panic("module replaced") }; module=d.Version } }
 if module != "v0.0.0-canonical" { panic("wrong artifact version: "+module) }
 failures:=0
 for _, op := range []string{"append","prepend","set"} {
  for _, foreignFirst := range []bool{false,true} {
   doc, err:=goquery.NewDocumentFromReader(strings.NewReader(string(input))); if err!=nil {panic(err)}
   ids:=[]string{"#html","#svg"}; if foreignFirst { ids=[]string{"#svg","#html"} }
   selected:=doc.Find(ids[0]).Add(ids[1]); expected:=""
   switch op {case "append":selected.AppendHtml(`<b>new</b>`);expected=`old<b>new</b>`;case "prepend":selected.PrependHtml(`<b>new</b>`);expected=`<b>new</b>old`;case "set":selected.SetHtml(`<b>new</b>`);expected=`<b>new</b>`}
   htmlB,svgB:=doc.Find("#html b").Length(),doc.Find("#svg b").Length()
   htmlText:=doc.Find("#html").Text();pass:=htmlB==0 && svgB==1 && htmlText==expected
   if !pass {failures++}
   encoded,_:=json.Marshal(map[string]interface{}{"operation":op,"foreignFirst":foreignFirst,"htmlB":htmlB,"svgB":svgB,"htmlText":htmlText,"expectedHTMLText":expected,"pass":pass,"runtime":runtime.Version(),"artifactVersion":module});fmt.Println(string(encoded))
  }
 }
 if failures>0 {fmt.Fprintf(os.Stderr,"%d of 6 real package consumer cases failed\n",failures);os.Exit(1)}
 fmt.Println("PASS: 6 real package consumer cases, 18 namespace/text checks; no replace directive")
}
