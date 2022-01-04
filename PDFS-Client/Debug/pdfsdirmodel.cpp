#include "pdfsdirmodel.h"

PDFSDirModel::PDFSDirModel()
{
    dirName.clear();
    createDate.clear();
    memset(shaCode,0x0,sizeof(shaCode));
    isRoot = haveFile = haveDir = sorted = false;
    fileList.clear();
    dirList.clear();
}

void PDFSDirModel::SetDirInfo(QString DirName, short YY, short MM, short DD, char* SHA)
{
    sorted=false;
    //set direcory
    dirName=DirName;

    //set friectory create date
    createDate.append(YY).append("\\").append(MM).append("\\").append(DD);

    //set SHA265Code
    for(int i=0;i<64;i++)
        shaCode[i]=SHA[i];
}

void PDFSDirModel::ChangeDate(short YY, short MM, short DD)
{
    createDate.clear();
    createDate.append(YY).append("\\").append(MM).append("\\").append(DD);
}

void PDFSDirModel::ChangeDirName(QString NewDirName)
{
    dirName=NewDirName;
}

void PDFSDirModel::SetSHACode(char* SHA)
{
    for(int i=0;i<64;i++)
        shaCode[i]=SHA[i];
}

RSC PDFSDirModel::AddDir(PDFSDirModel &NewDir)
{
    if(dirList.contains(NewDir.dirName))
        return run_DirExist;
    dirList[NewDir.dirName]=NewDir;
    return run_NoError;
}

RSC PDFSDirModel::AddFile(PDFSFileModel &NewFile)
{
    if(fileList.contains(NewFile.FileName()))
        return run_FileExist;
    fileList[NewFile.FileName()]=NewFile;
    return run_NoError;
}

void PDFSDirModel::DelDir(QString DirName)
{
    dirList.erase(dirList.find(DirName));
}

void PDFSDirModel::DelFile(QString FileName)
{
    fileList.erase((fileList.find(FileName)));
}

QString PDFSDirModel::DirName() { return dirName; }

QString PDFSDirModel::CreateDate() { return createDate; }
const char* PDFSDirModel::SHA() { return shaCode; }
