#include "pdfsfilesystemmodel.h"

PDFSFileSystemModel::PDFSFileSystemModel()
{
    root = new PDFSDirModel;
    now = root;
}
