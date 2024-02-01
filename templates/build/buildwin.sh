update-alternatives --set x86_64-w64-mingw32-gcc /usr/bin/x86_64-w64-mingw32-gcc-posix
update-alternatives --set x86_64-w64-mingw32-g++ /usr/bin/x86_64-w64-mingw32-g++-posix

make clean
make distclean

cd depends
make HOST=x86_64-w64-mingw32 -j$(nproc)
cd ..
./autogen.sh
./configure --prefix=`pwd`/depends/x86_64-w64-mingw32
make -j$(nproc)

cp src/{COIN_NAME-cli.exe,COIN_NAME-tx.exe,COIN_NAMEd.exe,qt/COIN_NAME-qt.exe} ~/tool/windist/