find . -type f -exec chmod +rx {} \;

cd depends
make HOST=x86_64-w64-mingw32 -j$(nproc)
make HOST=x86_64-pc-linux-gnu -j$(nproc)
cd ..