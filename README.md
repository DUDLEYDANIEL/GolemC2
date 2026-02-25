** structure of the project**
1. Server => manage tasks and results, communication to agents , managing multiple agent from a central server
2. Agent => a Go binary execting the tasks, beaconing for tasks and sending results 
3. cli => command line utility for sending tasks and viewing results 


** Functionlity features **
1. Communication Channels:
	Agents maybe use https for communication or connetion , how to implement stealthy mode of comms (mimicing normal traffic, reducind sus)
2. Tasking and Control:
	Server sends commands or taks, Agents executes them and agent should be beaconing for other tasks
4. Data exfilration:
	Agents should the result of the tasks(sensitive ), how to transfer the data securely
5. Evasion Techniques:
	should implement evasion of detection like randomizing the beacon intervals, domain fronting etc...

** TASKS **
* The tasks can be reconissance
* port scans and version detctection 
* mapping the internal network structure 
* lateral movements and Data exfiltration 

** comms **
1. HTTPS
2. mTLS

** encryptions **
1. mTLS
2. per agent asymmmetric keys (RSA, Ed25519)
3. storing and managing of secure keys 

