
import { createContext, useContext, useState, useEffect } from "react";
import { useWebSocket } from "@/components/WebSocketProvidor";

interface UserData {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
  Country: string;
}

interface UserDataContextType {
  userData: UserData;
  setUserData: React.Dispatch<React.SetStateAction<UserData>>;
}

const UserDataContext = createContext<UserDataContextType | null>(null);

export const useUserData = () => {
  const context = useContext(UserDataContext);
  if (!context) {
    throw new Error("useUserData must be used within a UserDataProvider");
  }
  return context;
};

export const UserDataProvider = ({ children }: { children: React.ReactNode }) => {
  const { socket, isReady } = useWebSocket();
  const [userData, setUserData] = useState<UserData>({
    ID: "",
    FirstName: "",
    LastName: "",
    Email: "",
    Country: "",
  });
  useEffect(()=>{
      
         if(socket && isReady){
               socket.send( "getUser" );

        socket.onMessage((msg)=>{
    const response =  JSON.parse(msg)
    if (response.userData) {
        
     setUserData(response.userData)
    }
    })}},[socket, isReady]) 
  return (
    <UserDataContext.Provider value={{ userData, setUserData }}>
      {children}
    </UserDataContext.Provider>
  );
};
