#include "pdfsfilemodel.h"

PDFSFileModel::PDFSFileModel()
{
    fileName.clear();
    fileType.clear();
    createDate.clear();
    memset(shaCode,0x0,sizeof(shaCode));
    fileSize = 0;
}

void PDFSFileModel::SetFileInfo(QString RawFileName, short YY, short MM, short DD, char* SHA, size_t FileSize)
{
    //set filename and filetype
    fileName.clear();
    fileType.clear();
    int pos=RawFileName.length()-1;
    while(pos>=0 && RawFileName[pos]!='.')
        pos--;

    if(pos)
    {
        for(int i=0;i<pos;i++)
            fileName.append(RawFileName[i]);
        for(int i=pos+1;i<RawFileName.length();i++)
            fileType.append(RawFileName[i]);
    }
    else
        fileName=RawFileName;

    //set file create date
    createDate.append(YY).append("\\").append(MM).append("\\").append(DD);

    //set SHACode
    for(int i=0;i<64;i++)
        shaCode[i]=SHA[i];

    //set filesize
    fileSize = FileSize;
}

void PDFSFileModel::ChangeDate(short YY, short MM, short DD)
{
    createDate.clear();
    createDate.append(YY).append("\\").append(MM).append("\\").append(DD);
}

void PDFSFileModel::ChangeFileName(QString NewFileName)
{
    fileName=NewFileName;
}

void PDFSFileModel::SetSHACode(char* SHA)
{
    for(int i=0;i<64;i++)
        shaCode[i]=SHA[i];
}

QString PDFSFileModel::FileName() { return fileName; }

QString PDFSFileModel::FileType() { return fileType; }

QString PDFSFileModel::CreateDate() { return createDate; }

const char* PDFSFileModel::SHA() { return shaCode; }

size_t PDFSFileModel::FileSize() { return fileSize; }
