make clean
make distclean

cd depends
make HOST=x86_64-pc-linux-gnu -j$(nproc)
cd ..
./autogen.sh
# ./configure --prefix=`pwd`/depends/x86_64-pc-linux-gnu --enable-debug
./configure --prefix=`pwd`/depends/x86_64-pc-linux-gnu
make -j$(nproc)

cp src/{COIN_NAME-cli,COIN_NAME-tx,COIN_NAMEd,qt/COIN_NAME-qt} ~/tool/unixdist/

