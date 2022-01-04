QT       += core gui
QT       += network

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

CONFIG += c++11

# You can make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    ligin.cpp \
    loginheader.cpp \
    logout.cpp \
    main.cpp \
    pdfs_client.cpp \
    pdfsdirmodel.cpp \
    pdfsfilemodel.cpp \
    pdfsfilesystemmodel.cpp \
    serverconnect.cpp

HEADERS += \
    login.h \
    loginheader.h \
    logout.h \
    pdfs_client.h \
    pdfsdirmodel.h \
    pdfsfilemodel.h \
    pdfsfilesystemmodel.h \
    protocol.h \
    serverconnect.h

FORMS += \
    pdfs_client.ui

# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target
