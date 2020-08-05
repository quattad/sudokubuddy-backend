cd ../targets/

cat targets_users.http | vegeta attack -rate 0 -max-workers 100 -duration 5s | tee results.bin | vegeta report