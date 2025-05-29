# This Apple script is created to start consistency hashing backend , 
# frontend part and scaling the deployment to 1 replica at starting to serve client
# Note that here path should be set according to your machine , where you cloned it.
# This is only the demo script , We can add more methods of automation.
# Also This scipt is using iTerm2 terminal.
# if You Do not have iTerm2 , then use commented script instead.
# This script should be changed according to Operating system.

#!/bin/bash
osascript <<EOF
tell application "iTerm"
  activate
  set newWindow to (create window with default profile)

  tell current session of newWindow
    write text "cd ~/Desktop/goserver/external-app && go run main.go"
  end tell

  -- split horizontally (right pane)
  tell current session of newWindow
    set rightSession to (split horizontally with default profile)
    tell rightSession
      write text "sleep 3; cd ~/Desktop/goserver/frontend && source venv/bin/activate && python server.py;"
    end tell

  -- split right pane horizontally again (bottom pane)
    set bottomSession to (split vertically with default profile)
    tell bottomSession
      write text "sleep 5 && open http://localhost:8888 && kubectl scale deployment go-app-server --replicas=1 -n demo"
    end tell
  end tell
end tell
EOF


# This script is for mac terminal.

# osascript <<EOF
# tell application "Terminal"
#   activate

#   do script "cd ~/Desktop/goserver/external-app; go run main.go"

#   delay 2

#   tell application "System Events" to keystroke "t" using {command down}
#   delay 1
#   do script "cd ~/Desktop/goserver/frontend; source venv/bin/activate; python server.py" in front window

#   delay 5

#   tell application "System Events" to keystroke "t" using {command down}
#   delay 1
#   do script "open http://localhost:8888; sleep 5; kubectl scale deployment go-app-server --replicas=1 -n demo" in front window
# end tell
# EOF