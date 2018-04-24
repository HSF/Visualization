
Bool_t CheckVolume(TGeoVolume *vol) 
{
    Int_t n = 0;
    Bool_t isany = kFALSE;
    while (n < vol->GetNdaughters()) {
        TGeoNode *node = vol->GetNode(n);
        Bool_t vis = CheckVolume(node->GetVolume());
        if (!vis) {
            vol->RemoveNode(node);
        } else {
            n++;
            isany = kTRUE;
        }
    }
    return isany ? kTRUE : vol->IsVisible();
}

void geomAtlas() {
   TGeoManager::Import("http://root.cern.ch/files/atlas.root");
   //gGeoManager->DefaultColors();
   gGeoManager->SetMaxVisNodes(10000);
   //gGeoManager->SetVisLevel(4);
   gGeoManager->GetVolume("ATLS")->Draw("ogl");
   
   // export only visible volumes
   CheckVolume(gGeoManager->GetVolume("ATLS"));
   
   // no any compression
   TBufferJSON::ExportToFile("atlas2.json", gGeoManager->GetVolume("ATLS"));
   
   // maximal compression, including final gzip
   TBufferJSON::ExportToFile("atlas2.json.gz", gGeoManager->GetVolume("ATLS"), "23");
   
   // just for comparasion store in the ROOT file
   TFile *f = TFile::Open("atlas2.root", "recreate");
   f->WriteObject(gGeoManager->GetVolume("ATLS"), "atlas");
   delete f;
}
