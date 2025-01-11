"use client";

import {
  Calendar,
  ChevronUp,
  Home,
  Inbox,
  Search,
  Settings,
  User2,
} from "lucide-react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";

import { ThemeChanger } from "@/components/ui/theme";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import { useNavigate } from "react-router";
import { toast } from "react-toastify";
import { useWebSocket } from "./WebSocketProvidor";
import { useEffect, useState } from "react";
const items = [
  {
    title: "Home",
    url: "#",
    icon: Home,
  },
  {
    title: "Inbox",
    url: "#",
    icon: Inbox,
  },
  {
    title: "Calendar",
    url: "#",
    icon: Calendar,
  },
  {
    title: "Search",
    url: "#",
    icon: Search,
  },
  {
    title: "Settings",
    url: "#",
    icon: Settings,
  },
];

export function AppSidebar() {
  const navigate = useNavigate();
   const {socket, isReady}= useWebSocket();

   const [userData, setUserData] = useState({
       FirstName: "",
      LastName:"",
      Email:"",
      ID:""
   });
  const handleLogout = async () => {
    try {
      const response = await fetch("http://localhost:3001/Logout", {
        method: "POST",
        headers: { "Content-Type": "application/json" },

        credentials: "include",
      });

      if (response.ok) {
        navigate("/login");
      } else {
        toast.error("Wrong credentials");
      }
    } catch (error) {
      toast.error("Login Failed ");
    }
  };

 

  useEffect(()=>{
      
         if(socket && isReady){
               socket.send( "getUser" );

        socket.onMessage((msg)=>{
    const response =  JSON.parse(msg)
    if (response.userData) {
        
     setUserData(response.userData)
    }
    })}},[socket, isReady]) 
useEffect(() => {
}, [userData]); 
    
    return (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
          <ThemeChanger />
          <SidebarGroupLabel>Application</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton>
                  <User2 /> {userData.FirstName}
                  <ChevronUp className="ml-auto" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                side="top"
                className="w-[--radix-popper-anchor-width]"
              >
                <DropdownMenuItem>
                  <span>Account</span>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <span>Billing</span>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <span onClick={handleLogout}>Sign out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}

