#include "user_class.h"

#include "TBufferJSON.h"

void createJSON()
{
   UserClass obj;
   obj.fInt = 12;
   for (int n=0;n<10;++n)
     obj.fVect.emplace_back(n*0.5);
   obj.fStr = "AnyString";
   TString json = TBufferJSON::ToJSON(&obj);

   printf("json \n%s\n", json.Data());


}
