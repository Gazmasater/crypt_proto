git rm --cached cmd/arb/metrics/arb_metrics.csv

echo "cmd/arb/metrics/*.csv" >> .gitignore

git add .gitignore
git commit --amend --no-edit


git push origin new_arh --force



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ git rm --cached cmd/arb/metrics/arb_metrics.csv
fatal: pathspec 'cmd/arb/metrics/arb_metrics.csv' did not match any files
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ 




