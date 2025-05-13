#RUN FOR WINDOWS
```
CGO_ENABLED=1 GOOS=windows CC=x86_64-w64-mingw32-gcc go run -tags nolibopusfile . 
```

//   yt-dlp -f "bestaudio[ext=webm][acodec=opus]" -g

# Notes
- Move to this https://github.com/gotranspile/opus
