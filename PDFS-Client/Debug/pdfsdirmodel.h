#ifndef PDFSDIRMODEL_H
#define PDFSDIRMODEL_H

#include <QMap>
#include <algorithm>

#include "pdfsfilemodel.h"
#include "protocol.h"

class PDFSDirModel
{
private:
    QString dirName;
    QString createDate;
    char shaCode[64];

    bool isRoot;
    bool haveFile;
    bool haveDir;
    bool sorted;

    QMap<QString,PDFSFileModel> fileList;
    QMap<QString,PDFSDirModel> dirList;

public:
    PDFSDirModel();
    void SetDirInfo(QString DirName, short YY, short MM, short DD, char* SHA);
    void ChangeDate(short YY, short MM, short DD);
    void ChangeDirName(QString NewFileName);
    void SetSHACode(char* SHA);
    RSC AddDir(PDFSDirModel &NewDir);
    RSC AddFile(PDFSFileModel &NewFile);
    void DelDir(QString DirName);
    void DelFile(QString FileName);

public:
    QString DirName();
    QString CreateDate();
    const char* SHA();

public:
    bool friend operator < (PDFSDirModel A,PDFSDirModel B)
    {
        return A.dirName<B.dirName;
    }
};

#endif // PDFSDirMODEL_H
