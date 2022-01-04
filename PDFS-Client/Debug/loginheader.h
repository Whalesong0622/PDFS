#ifndef LOGINHEADER_H
#define LOGINHEADER_H

#include <QString>
#include <QBitArray>

#include "protocol.h"

class LoginHeader
{
public:
    LoginHeader();

private:
    char op;
    char userNameLength;
    QString userName;
    QString password;

    QByteArray header;
    QString rawHeader;

    bool headerReady;
public:
    void set_op(int Op);
    void set_UserName(QString UserName);
    void set_Password( QString Password);
    QByteArray get_Header();
    QString get_RawHeader();

private:
    void rebuildHeader();
};

#endif // LOGINHEADER_H
