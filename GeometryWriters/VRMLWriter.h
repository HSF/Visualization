/**************************************************************************
 * BASF2 (Belle Analysis Framework 2)                                     *
 * Copyright(C) 2016 - Belle II Collaboration                             *
 *                                                                        *
 * Author: Leo Piilonen                                                   *
 * Contributors: Leo Piilonen                                             *
 *                                                                        *
 * This software is provided "as is" without any warranty.                *
 **************************************************************************/

#ifndef VRMLWRITER_H
#define VRMLWRITER_H

#include <string>
#include <fstream>

#include "G4Transform3D.hh"  // Using 'class G4Transform3D' conflicts with a typedef in here
class G4VPhysicalVolume;
class G4LogicalVolume;
class G4VSolid;
class G4AffineTransform;
class G4Polyhedron;
class HepPolyhedron;

/** The VRML-writer module.
 *
 * This module goes through all volumes of the constructed GEANT4
 * geometry and writes a VRML 2.0 file.  Unlike the VRML writer in
 * GEANT4, which creates a flat model, this writes the hierarchical model.
 *
 * Prerequisite: This module requires a valid GEANT4 geometry.
 *
 */
class VRMLWriter {

public:

  //! Constructor [empty]
  VRMLWriter() {}

  //! Destructor [empty]
  ~VRMLWriter() {}

  //! Write the geometry in the VRML 2.0 (text) format
  //! @param outputFilename: user-specified output filename [empty string -> "geometry.wrl"]
  //! @return true if the VRML file was written; false otherwise
  bool doit(std::string outputFilename);

private:

  //! Create unique and legal name for each solid
  void assignName(std::vector<std::string>*, unsigned int, const G4String&, int);

  //! Emit VRML for each solid
  void describeSolid(G4VSolid*, const std::string&, bool);

  //! Emit VRML for each logical volume
  void describeLogicalVolume(G4LogicalVolume*, const std::string&, const std::string&, bool);

  //! Access next physical volume in the tree (recursive)
  void describePhysicalVolume(G4VPhysicalVolume*);

  //! Emit VRML for a physical volume (recursive)
  void descendAndDescribe(G4VPhysicalVolume*, const std::string&, int);

  //! Emit VRML for the start of the file
  void writePreamble(void);

  //! Emit VRML for the solid's polyhedron
  void writePolyhedron(const G4Polyhedron*, const std::string&);

  //! Emit VRML for each daughter of a logical volume
  void writePhysicalVolume(const G4VPhysicalVolume*, const std::string&, const std::string&, bool);

  //! Create polyhedron for a boolean solid (recursive)
  HepPolyhedron* getBooleanSolidPolyhedron(G4VSolid*);

  //! Output file
  std::ofstream m_File;

  //! Modified (legal-character and unique) physical-volume name
  std::vector<std::string>* m_PVName;

  //! Modified (legal-character and unique) logical-volume name
  std::vector<std::string>* m_LVName;

  //! Modified (legal-character and unique) solid name
  std::vector<std::string>* m_SolidName;

  //! Flag to indicate that a solid can be rendered as a VMRL cylinder
  std::vector<bool>* m_IsCylinder;

  //! Indices (in G4PhysicalVolumeStore) of the logical volume's physical-volume daughters
  std::vector<std::vector<int> >* m_PVIndex;

  //! Flag to indicate that the logical volume has already been written
  std::vector<bool>* m_LVWritten;

  //! Flag to indicate that the physical volume has already been written
  std::vector<bool>* m_PVWritten;

};

#endif
