import java.io.*;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.UnknownHostException;

public class Main {

    public static void main(String[] args) throws Exception {
            final int myPort = 10001;
            ServerSocket ssock = new ServerSocket(myPort);
            System.out.println("Listening on port " + myPort );

        while (true) {
            Socket sock = ssock.accept();
            System.out.println("Someone has made socket connection");
            OneConnection client = new OneConnection(sock);
            String s = client.getRequest();
        }
    }

}

class OneConnection {
    Socket sock;
    BufferedReader in = null;
    DataOutputStream out = null;

    OneConnection(Socket sock) throws Exception {
        this.sock = sock;
        in = new BufferedReader(new InputStreamReader(sock.getInputStream()));
        out = new DataOutputStream(sock.getOutputStream());
    }

    String getRequest() throws Exception {
        String s = null;
        while ((s = in.readLine()) != null) {
            System.out.println("got: " + s);
        }
        return s;
    }
}

