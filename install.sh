#!/bin/bash
os=$(uname)
arch=$(uname -m)

if [ "$os" != "Darwin" -a "$os" != "Linux" ]
then
    echo "Unsupported operating system $os"    
    exit 1
fi

if [ "$arch" == "arm64" -o "$arch" == "aarch64" -o "$arch" == "armv8b" -o "$arch" == "armv8l" ]
then
    arch="arm64"
elif [ "$arch" == "i386" -o "$arch" == "i686" ]
then
    arch="386"
elif [ "$arch" == "x86_64" ]
then
    arch="amd64"
else
    echo "Unsupported architecture $arch"
    exit 1
fi

os=$(echo $os | awk '{print tolower($0)}')
download_tag="donggu-$os-$arch.tar"

echo "Downloading for OS: $os, Arch: $arch"

download_url=$(curl -s https://api.github.com/repos/maasasia/donggu/releases/latest | grep "/$download_tag" | cut -f 2- -d : | cut -f 2 -d \")
echo "Downloading from $download_url"

rm -f donggu.tar
wget -O donggu.tar "$download_url" || exit 1

echo "Download complete."
rm -rf donggu
mkdir donggu
tar xvf donggu.tar -C donggu || exit 1
rm donggu.tar

echo ""
echo "Donggu installed at $(pwd)/donggu"
echo 'Either add this directory to $PATH, or move it to somewhere included in $PATH.'

