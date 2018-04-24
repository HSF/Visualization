#include "user_class.h"

#include "TBufferJSON.h"
#include <fstream>
#include <cstdio>

void createJSON()
{
   // create and fill object
   UserClass obj;
   obj.fInt = 12;
   for (int n=0;n<10;++n)
     obj.fVect.emplace_back(n*0.5);
   obj.fStr = "AnyString";

   // create JSON representation for the object
   TString json = TBufferJSON::ToJSON(&obj);

   // print to std output
   printf("json:\n%s\n", json.Data());

   // save to the file
   std::ofstream ofs("file.json");
   ofs << json;
   ofs.close();

}
