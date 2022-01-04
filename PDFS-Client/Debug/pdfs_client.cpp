#include "pdfs_client.h"
#include "ui_pdfs_client.h"

PDFS_Client::PDFS_Client(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::PDFS_Client)
{
    ui->setupUi(this);
    ui->ServerState->setEnabled(false);
    ui->ServerState->setText("disconnect");

    ui->FileTree->setColumnCount(3);
    //ui->FileTree->setColumnHidden(0,true);
    QList<QString> headerName;
    headerName.append("FileName");
    headerName.append("Size");
    headerName.append("Create data");
    ui->FileTree->setHeaderLabels(headerName);
    QList<QTreeWidgetItem *> items;
    for (int i = 0; i < 10; ++i)
    {
        QStringList sl;
        sl.append(QString("item: %1").arg(i));
        sl.append(QString("item: %1").arg(i*10));
        items.append(new QTreeWidgetItem(static_cast<QTreeWidget *>(nullptr), sl));
    }
    ui->FileTree->insertTopLevelItems(0,items);
}

PDFS_Client::~PDFS_Client()
{
    delete ui;
}

void PDFS_Client::on_ServerConnect_clicked()
{
    ui->DebugInfo->append("Try to connect server:");
    ServerConnect server;
    ui->DebugInfo->append(server.ConnectTest("43.132.181.175"));
}

void PDFS_Client::on_Login_clicked()
{
    tcpHeader.set_op(0);
    tcpHeader.set_UserName(ui->UserName->text());
    tcpHeader.set_Password(ui->Password->text());
    ui->DebugInfo->append("Username:    " + ui->UserName->text());
    ui->DebugInfo->append("Password:    " + ui->Password->text());
    ui->DebugInfo->append("Raw Header:\n" + tcpHeader.get_RawHeader()+"\n");
    ui->DebugInfo->append("Hex Header:\n" + tcpHeader.get_Header().toHex()+"\n");
}

void PDFS_Client::on_FileTree_itemClicked(QTreeWidgetItem *item, int column)
{
    ui->DebugInfo->append("user clicked dir:    " + item->data(0,0).toString());
}
