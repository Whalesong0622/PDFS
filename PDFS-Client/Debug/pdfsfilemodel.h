#ifndef PDFSFILEMODEL_H
#define PDFSFILEMODEL_H

#include <QString>
#include <QList>

class PDFSFileModel
{
private:
    QString fileName;
    QString fileType;
    QString createDate;
    char shaCode[64];
    size_t fileSize;

public:
    PDFSFileModel();
    void SetFileInfo(QString RawFileName, short YY, short MM, short DD, char* SHA, size_t FileSize);
    void ChangeDate(short YY, short MM, short DD);
    void ChangeFileName(QString NewFileName);
    void SetSHACode(char* SHA);

public:
    QString FileName();
    QString FileType();
    QString CreateDate();
    const char* SHA();
    size_t FileSize();

public:
    bool friend operator < (PDFSFileModel A,PDFSFileModel B)
    {
        return A.fileName<B.fileName;
    }
};

#endif // PDFSFILEMODEL_H
