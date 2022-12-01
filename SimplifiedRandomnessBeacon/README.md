# Randomness-Beacon
VRF,PBFT,Time Commitment

## Go-Algorand
### Initial environment setup:

```
git clone https://github.com/algorand/go-algorand
cd go-algorand
sudo ./scripts/configure_dev.sh
sudo ./scripts/buildtools/install_buildtools.sh
```

### build

```
sudo make install
```

### Note
Do not forget to change the local Go-Algorand path in `ConsensusNode/go.mod` and `EntropyProvider/go.mod`

## Randomness Beacon
```
sudo git clone https://github.com/PsyduckLiu/Randomness-Beacon.git
cd Randomness-Beacon
cd Command
```

### start randomness beacon
```
 ./start_RandBeacon.sh
```
### stop randomness beacon
```
 ./end_RandBeacon.sh
```

## Notes

If you want to modify the number of Entropy Providers, just change the settings in `Command/start_RandBeacon.sh` and `Command/end_RandBeacon.sh`.

If you want to modify the number of Consensus Nodes, please change the settings in `Command/start_RandBeacon.sh` and `Command/end_RandBeacon.sh`, and the parameter `MaxFaultyNode` in `ConsensusNode/util/util.go` and `EntropyProvider/util/util.go`.