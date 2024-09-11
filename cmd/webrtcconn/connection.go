package webrtcconn

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func CreatePeerConnection() *webrtc.PeerConnection {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatalf("Failed to create peer connection: %v", err)
	}

	return pc
}

func HandleSignaling(pc *webrtc.PeerConnection, roomID string) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatalf("Failed to connect to signaling server: %v", err)
	}
	defer conn.Close()

	conn.WriteJSON(map[string]string{"type": "join", "room": roomID})

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println(err)
			break
		}

		switch msg["type"] {
		case "offer":
			offer := msg["sdp"].(string)
			sessionDescription := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offer}
			pc.SetRemoteDescription(sessionDescription)
			answer, err := pc.CreateAnswer(nil)
			if err != nil {
				log.Fatalf("Failed to create answer: %v", err)
			}
			pc.SetLocalDescription(answer)
			conn.WriteJSON(map[string]interface{}{
				"type": "answer",
				"sdp":  answer.SDP,
			})

		case "answer":
			answer := msg["sdp"].(string)
			sessionDescription := webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer}
			pc.SetRemoteDescription(sessionDescription)

		case "ice-candidate":
			candidate := msg["candidate"].(string)
			pc.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate})

		default:
			fmt.Println("Unknown signaling message type:", msg["type"])
		}
	}
}
