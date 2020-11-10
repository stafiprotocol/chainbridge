#Instructions of starting bridge.

**Prepare** _This only needs to be done once_
1. Generate keystore file. There are several mini steps.
    1. `make install`
    2. `chainbridge account gensub` to generate keystore for substrate, mnemonics or rawseed is needed.
    3. `chainbridge account geneth` to generate keystore for ethereum, private key is needed.
    4. password is also needed for last two mini steps, password for each keystore file could be different.
2. Run setup test to finish contract configuration. 
    `go test -v chains/ethereum/setup_test.go`, the test file must be modified before running this command, there are two places to fill.
    1. keystorePath.  line 29. 
    2. password. line 34. Should be the same as the one you just used to generate the keystore files.
    
**Start**
3. Start Bridge:
`chainbridge --config configfile`, once again password should be given for the keystore file.
    