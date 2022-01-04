#ifndef PDFS_CLIENT_H
#define PDFS_CLIENT_H

#include <QMainWindow>
#include <QTcpSocket>
#include <QFileSystemModel>
#include <QTreeWidgetItem>
#include <QDir>
#include <QList>

#include "serverconnect.h"
#include "loginheader.h"

QT_BEGIN_NAMESPACE
namespace Ui { class PDFS_Client; }
QT_END_NAMESPACE

class PDFS_Client : public QMainWindow
{
    Q_OBJECT

private:
    QTcpSocket *Client;
    LoginHeader tcpHeader;
    QFileSystemModel *model;

public:
    PDFS_Client(QWidget *parent = nullptr);
    ~PDFS_Client();

private slots:
    void on_ServerConnect_clicked();
    void on_Login_clicked();

    void on_FileTree_itemClicked(QTreeWidgetItem *item, int column);

private:
    Ui::PDFS_Client *ui;
};
#endif // PDFS_CLIENT_H
