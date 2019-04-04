 #!/bin/bash
        for i in `seq 4010 4015`;
        do
                nohup ./bin/sbottle-linux-amd64 --port $i &>/dev/null &
        done    
