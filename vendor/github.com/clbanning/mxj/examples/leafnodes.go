// https://groups.google.com/forum/#!topic/golang-nuts/pj0C5IrZk4I

package main

import (
   "fmt"

   "github.com/clbanning/mxj"
)

func main() {
   j := `{"jsonData":{
      "DataReference":[
         {
            "ParameterType":"test",
            "Applicationtype":[
               {
                  "Application1":{
                     "ApplicationName":"app1",
                     "Param1":{
                        "Name":"app1.param1"
                     },
                     "Param2":{
                        "Name":"app1.param2"
                     }
                  },
                  "Application2":{
                     "ApplicationName":"app2",
                     "Param1":{
                        "Name":"app2.param1"
                     },
                     "Param2":{
                        "Name":"app2.param2"
                     }
                  }
               }
            ]
         }
      ]
   }}`

   // unmarshal into a map
   m, err := mxj.NewMapJson([]byte(j))
   if err != nil {
      fmt.Println("err:", err)
      return
   }
   mxj.LeafUseDotNotation()
   l := m.LeafNodes()
   for _, v := range l {
      fmt.Println("path:", v.Path, "value:", v.Value)
   }
   /*
      Output (sequence not guaranteed):
      path: jsonData.DataReference.0.ParameterType value: test
      path: jsonData.DataReference.0.Applicationtype.0.Application1.ApplicationName value: app1
      path: jsonData.DataReference.0.Applicationtype.0.Application1.Param1.Name value: app1.param1
      path: jsonData.DataReference.0.Applicationtype.0.Application1.Param2.Name value: app1.param2
      path: jsonData.DataReference.0.Applicationtype.0.Application2.ApplicationName value: app2
      path: jsonData.DataReference.0.Applicationtype.0.Application2.Param1.Name value: app2.param1
      path: jsonData.DataReference.0.Applicationtype.0.Application2.Param2.Name value: app2.param2
   */
}
