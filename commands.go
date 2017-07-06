package main

func (lp *LocalPeer) Ping(address *Address) HTTPResponse {

}

func (lp *LocalPeer) Announce(address *Address) HTTPResponse {

}

func (lp *LocalPeer) Search(query string, address *Address, page int, recursive bool) HTTPResponse {

}

func (lp *LocalPeer) Recent(address *Address, page int) HTTPResponse {

}

func (lp *LocalPeer) Popular(address *Address, page int) HTTPResponse {

}

func (lp *LocalPeer) Mirror(address *Address) HTTPResponse {

}

func (lp *LocalPeer) MirrorProgress(address *Address) HTTPResponse {

}

func (lp *LocalPeer) Index(address *Address, since int) HTTPResponse {

}

func (lp *LocalPeer) AddPost(post *Post, index bool) HTTPResponse {

}

func (lp *LocalPeer) Resolve(address *Address) HTTPResponse {

}

func (lp *LocalPeer) Bootstrap(address *Address) HTTPResponse {

}

func (lp *LocalPeer) Peers() HTTPResponse {

}

func (lp *LocalPeer) AddPeer(address *Address) HTTPResponse {

}

func (lp *LocalPeer) Set(identity Identity) HTTPResponse {

}

func (lp *LocalPeer) Get(key string) HTTPResponse {

}

func (lp *LocalPeer) Explore() HTTPResponse {

}

func (lp *LocalPeer) Map() HTTPResponse {

}
