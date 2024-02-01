import time
from charpreffix import preffix_dict

stamp = int(time.time())
print(f"getting timestamp: {stamp}")



chara = input("find char_prefix Bitcoin: ")

for preffix in preffix_dict:
    chardict = preffix["symbol"]
    chardict = chardict.split(",")
    
    datalist = []
    for dd in chardict:
        if dd.find("-") != -1:
            dds = dd.split("-")
            dds = list(map(lambda x:ord(x), dds))
            
            for cint in range(dds[0],dds[0]+1):
                datalist.append(chr(cint))
            
        else:
            datalist.append(dd)
    
    
    try:
        datalist.index(chara)
        print(preffix["decimal"])
    except ValueError:
        
        continue
