from Crypto.PublicKey import RSA
from Crypto.Hash import SHA512
import json
from collections import OrderedDict
from Crypto.Signature import pkcs1_15

def sign(sk, msg):
    '''
    sign string msg with corresponding secret key (RsaKey object) in pycryptodome library
    return hex value of signature with RSA pkcs1 v1.5
    '''
    msgToUTF = msg.encode("UTF-8")
#     print(msgToUTF)
    msgHash = SHA512.new(msgToUTF)
#     print(msgHash)
    sig = pkcs1_15.new(sk).sign(msgHash).hex()
#     print(sig)
    return sig


def verify(pk, msg, sig):
    '''
    check sign is made for msg using public key PK, string MSG, 
    and byte string SIGN.
    suppose publicExponent is fixed at 0x10001.
    return boolean
    '''
    msgHash = SHA512.new(msg.encode("UTF-8"))
    try:
        pkcs1_15.new(pk).verify(msgHash, bytes.fromhex(sig))
        return True
    except (ValueError, TypeError):
        return False

def load_secret_key(fname):
    '''
    load json information of secret key from fname. 
    This returns RSA key object of pycrytodome library.
    '''
    with open(fname, 'r') as f:
        file = json.load(f)
    modulus = int(file["modulus"], 16)
    publicExponent = int(file["publicExponent"], 16)
    privateExponent = int(file["privateExponent"], 16)
    rsaComponent = (modulus, publicExponent, privateExponent)
    sk = RSA.construct(rsaComponent, True)
    return sk


def create_secret_key(fname):
    '''
    Create a secret key: [hint] RSA.generate().
    Save the secret key in json to a file named "fname". 
    This returns RSA key object of pycrytodome library.
    '''
    publicExponent = 0x10001
    keyLength = 2048
    sk = RSA.generate(keyLength, e = publicExponent)
    pk = sk.publickey()
    file = OrderedDict()
    file["modulus"] = sk.n
    file["publicExponent"] = sk.e
    file["privateExponent"] = sk.d
    with open(fname, 'w') as f:
        json.dump(file, f, ensure_ascii=False, indent="t/")
    return sk


def get_hash(msg):
    '''
    return hash hexdigest for string msg with 0x. ex) 0x1a2b...
    '''
    msgHash = SHA512.new(msg.encode("UTF-8"))
    return msgHash.hexdigest()


def get_pk(sk):
    '''
    return pk using modulus of given RsaKey object sk.
    '''
    pk = sk.publickey()
    return pk
