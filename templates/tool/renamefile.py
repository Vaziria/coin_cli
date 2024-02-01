import os
import json

replacestr = ["tonnage", "vishai"]
replaceWord = [
    ["PDC-Repository", "vishcoin"],
    ["TonnageOfficial", "vishcoin"],
    ["Vaziria", "vishcoin"],
    ["Tonnage", "Vishai"],
    ["TonnageCrashInfo", "VishaiCrashInfo"],
    ["tonnage", "vishai"],
    ["TONNAGE", "VISHAI"],
    ["TNN", "VISH"],
    ["swampschool.org", "vishcoin.com"],
    ["vishaipro.tech", "vishcoin.com"]
]

RENAME_DICT_NAME = "rename_dictionary.txt"

ignores = [
    ".git",
    "depends",
    "depends/work",
    "depends/x86_64-w64-mingw32",
    ".git",
    "depends/x86_64-pc-linux-gnu",
    "src/.libs",
    "src/.deps",
    "/depends/sources",
]

def save_replace_word(file, word):
    with open(RENAME_DICT_NAME, "a+") as out:
        data = {
            "type": "replace_word",
            "file": file,
            "word": word,
        }
        
        text = json.dumps(data)
        
        out.write(text)
        out.write("\n")
        

def check_contain_word(file):
    with open(file, "r", encoding='utf-8') as out:
        try:
            data = out.read()
        except Exception as e:
            print(e)
            
            return (False, "")
        
    found = False
    for word in replaceWord:
        wfound = data.find(word[0]) != -1
        found = wfound or found
        data = data.replace(word[0], word[1])
        if wfound:
            save_replace_word(file, word)
        
    
    return (found, data)
    




def ignored(text):
    text = text.replace("\\", "/")
    for c in ignores:
        if text.find(c) != -1:
            return True
        
    return False



def check_path_contain(text):
    
    if text.find(replacestr[0]) != -1:
        return True
        
    return False



def save_dict_replace(oldfile, newfile):
    with open(RENAME_DICT_NAME, "a+") as out:
        data = {
            "type": "rename",
            "old_file": oldfile,
            "new_file": newfile
        }
        
        text = json.dumps(data)
        
        out.write(text)
        out.write("\n")
        


        

for root, dirs, files in os.walk("."):
    for dirt in dirs:
        dirt = os.path.join(root, dirt)
        
        if ignored(dirt):
            continue
        
        if check_path_contain(dirt):
            print(dirt)
            dstdir = dirt.replace(replacestr[0], replacestr[1])
            
            os.rename(dirt, dstdir)
    
    for file in files:
        file = os.path.join(root, file)
        
        
        if ignored(file):
            continue
        
        found, data = check_contain_word(file)
        if found:
            print("replacing", file)
            with open(file, "w+") as out:
                out.write(data)
            
            
        
        if check_path_contain(file):
            print(file)
            dstfile = file.replace(replacestr[0], replacestr[1])
            
            os.rename(file, dstfile)
            save_dict_replace(file, dstfile)
            
        