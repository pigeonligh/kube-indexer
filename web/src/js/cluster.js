export function GetCluster(){
    let path = window.location.pathname;
    while (path.startsWith('/')) {
        path = path.slice(1);
    }
    if (path === "") {
        return "default"
    }
    return path
};

export function ChangeCluster(cluster) {
    window.location.assign("/"+cluster)
}