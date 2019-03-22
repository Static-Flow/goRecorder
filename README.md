# goRecorder
During pentesting I often miss screenshots of events for reports due to the quick pace of testing and a lack of foreknowledge about what will be important. To remedy that problem (and also to teach myself go) I built a command line tool that implements the "clip that" functionality of gaming consoles to allow me to save the last minute of screen activity as images to later view.

By default the tool stores an image a second for each screen in a 60 second sliding window, which can both be changed. Once the tool is running you can type F1 with the tool window in focus to save the current buffer to disk for later viewing. If you step away from testing, or don't want the screen recorded, type F2 to pause collection and type F2 again to resume. 
To exit, simply type ESC.
