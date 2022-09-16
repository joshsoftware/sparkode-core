

-- Run time

isolate --cg -b 1 --init

cd /var/local/lib/isolate/1
touch stdin.txt
touch stdout.txt
touch stderr.txt
cd box
echo "puts 'hello'" >> script.rb
echo "/usr/local/ruby-2.7.0/bin/ruby script.rb" >> run

cd ~

isolate --cg -s -b 1 -M /var/local/lib/isolate/1/metadata.txt -t 5.0 -x 1.0 -w 10.0 -k 64000 -p60 --cg-timing --cg-mem=128000 -f 1024 -E HOME=/tmp -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -E LANG -E LANGUAGE -E LC_ALL -E JUDGE0_HOMEPAGE -E JUDGE0_SOURCE_CODE -E JUDGE0_MAINTAINER -E JUDGE0_VERSION -d /etc:noexec --run -- /bin/bash run < /var/local/lib/isolate/1/stdin.txt > /var/local/lib/isolate/1/stdout.txt 2> /var/local/lib/isolate/1/stderr.txt


-- Compile time
{
 source_file: "main.go",
 compile_cmd: "GOCACHE=/tmp/.cache/go-build /usr/local/go-1.13.5/bin/go build %s main.go",
 run_cmd: "./main"
}

isolate --cg -b 5 --init
cd /var/local/lib/isolate/5
touch stdin.txt
touch stdout.txt
touch stderr.txt
touch compile_output.txt 
cd box

echo "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Print(\"hello\");\n\n}" >> main.go
echo "GOCACHE=/tmp/.cache/go-build /usr/local/go-1.19.1/bin/go build main.go" >> compile

isolate --cg -s -b 5 -M /var/local/lib/isolate/5/metadata.txt --stderr-to-stdout -i /dev/null -t 15.0 -x 0 -w 20.0 -k 128000 -p120 --cg-timing --cg-mem=512000 -f 4096 -E HOME=/tmp -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -E LANG -E LANGUAGE -E LC_ALL -E JUDGE0_HOMEPAGE -E JUDGE0_SOURCE_CODE -E JUDGE0_MAINTAINER -E JUDGE0_VERSION -d /etc:noexec --run -- /bin/bash compile > /var/local/lib/isolate/5/compile_output.txt 

echo "./main" >> run

isolate --cg -s -b 5 -M /var/local/lib/isolate/5/metadata.txt -t 5.0 -x 1.0 -w 10.0 -k 64000 -p60 --cg-timing --cg-mem=128000 -f 1024 -E HOME=/tmp -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -E LANG -E LANGUAGE -E LC_ALL -E JUDGE0_HOMEPAGE -E JUDGE0_SOURCE_CODE -E JUDGE0_MAINTAINER -E JUDGE0_VERSION -d /etc:noexec --run -- /bin/bash run < /var/local/lib/isolate/5/stdin.txt > /var/local/lib/isolate/5/stdout.txt 2> /var/local/lib/isolate/5/stderr.txt 


