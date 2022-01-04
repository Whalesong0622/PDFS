#include "loginheader.h"

LoginHeader::LoginHeader()
{
    this->op = 0;
    this->userNameLength = 0;
    this->header.clear();
    this->password.clear();
    this->userName.clear();
    this->headerReady = false;
}

void LoginHeader::set_op(int Op)
{
    op=Op;
    headerReady = false;
}

void LoginHeader::set_UserName(QString UserName)
{
    userName=UserName;
    userNameLength=UserName.length();
    headerReady = false;
}

void LoginHeader::set_Password(QString Password)
{
    password=Password;
    headerReady = false;
}

QByteArray LoginHeader::get_Header()
{
    if(!headerReady)
    {
        rebuildHeader();
        headerReady=true;
    }
    return header;
}

QString LoginHeader::get_RawHeader()
{
    if(!headerReady)
    {
        rebuildHeader();
        headerReady=true;
    }
    return rawHeader;
}

void LoginHeader::rebuildHeader()
{
    rawHeader.clear();
    rawHeader.append(op);
    rawHeader.append(userNameLength);
    rawHeader.append(userName);
    rawHeader.append(password);
    header = rawHeader.toLatin1();
}
