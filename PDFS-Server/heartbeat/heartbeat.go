package heartbeat

import "time"

func HeartBeatTimer(){
	for{
		HeartBeat()
		time.Sleep(time.Second*20)
	}
}

func HeartBeat(){

}