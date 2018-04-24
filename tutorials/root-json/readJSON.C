#include "user_class.h"

#include "TBufferJSON.h"

#include <cstdio>
#include <string>
#include <fstream>
#include <streambuf>

void readJSON()
{
   // read text file
   std::ifstream ifs("file.json");
   std::string json((std::istreambuf_iterator<char>(ifs)),
                    std::istreambuf_iterator<char>());

   // print to std output
   printf("json:\n%s\n", json.c_str());

   // declare empty pointer on objeect
   UserClass *obj = nullptr;

   // create object using JSON declaration
   TBufferJSON::FromJSON(obj, json.c_str());

   if (!obj) {
      printf("Fail to read object from file.json\n");
   } else {

      // print to std output
      printf("obj->fInt = %d\n", obj->fInt);
      printf("obj->fVect.size() = %u\n", (unsigned) obj->fVect.size());
      for (unsigned n=0; n<obj->fVect.size(); ++n)
         printf("   [%u] = %f\n", n, obj->fVect[n]);
      printf("obj->fStr = \"%s\"\n", obj->fStr.c_str());

      // not to forget delete object
      delete obj;
   }


}
