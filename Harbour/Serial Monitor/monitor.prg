#include "hbcom.ch"
#include "inkey.ch"

FUNCTION Main( )

   LOCAL cString := ""
   LOCAL nTimeOut := 500
   LOCAL nResult
   
   local line:=""
   LOCAL nPort := 0
   local pa:=DetectSerialPorts()
   local lc:=0
   local ik:=0
   local cp:=0
   local gpsport:=""
   local rcport:=""
   local senport:=""

   ?os()
   if len(pa)>0 

      DO WHILE cp < len(pa)
         cp++
         lc=0
         nport=pa[cp]
         IF ! hb_comOpen( nPort )
            ? "Cannot open port:", nPort, hb_comGetDevice( nPort ), ;
               "error: " + hb_ntos( hb_comGetError( nPort ) )
         ELSE
            ? "port:", hb_comGetDevice( nPort ), "opened"
            IF ! hb_comInit( nPort, 9600, "N", 8, 1 )
               ? "Cannot initialize port to: 9600:N:8:1", ;
                  "error: " + hb_ntos( hb_comGetError( nPort ) )
            ELSE
         
               do while .t.
                  ik=inkey()
                  if ik = 27 
                     exit
                  endif       
                  cString := Space( 1 )
                  nTimeOut := 500 
                  nResult := hb_comRecv( nPort, @cString, hb_BLen( cString ),nTimeOut )
                  IF nResult == 1
                     if asc(cstring)=13
                        lc++
                        do case
                           case substr(line,1,3)="$GP"
                              gpsport=str(nport)
                              hb_comClose( nPort )
                              line=""
                              exit
                           case substr(line,1,3)="CH1"
                              rcport=str(nport)
                              hb_comClose( nPort )
                              line=""
                              exit
                           case substr(line,1,3)="DIS"
                              senport=str(nport)
                              hb_comClose( nPort )
                              line=""
                              exit
                        endcase
                        line=""
                        if lc>100
                           hb_comClose( nPort )
                           exit
                        endif
                     else
                        line=line+cstring   
                     endif

                  ENDIF
               enddo
            ENDIF
            hb_comClose( nPort )
         ENDIF
      enddo
   else
      ?"No Serial Ports Found"
   endif
   if len(gpsport)=0
      ?"GPS port not found"
   else
      ?"GPS Port :",gpsport
   endif
   if len(rcport)=0
      ?"RC port not found"
   else
      ?"RC Port :",rcport
   endif
   if len(senport)=0
      ?"Sensor port not found"
   else
      ?"Sensor Port :",senport
   endif
RETURN

FUNCTION DetectSerialPorts()
   local pa:=array(0)
   local x:=25
   local nPort
   local cPortName:="/dev/ttyACM"
   local cPort:=""
   DO WHILE x > 0
      nPort=x
      IF hb_comOpen( nPort )
         aadd(pa,x)
         hb_comClose( nPort )
      ELSE
         cPort=cPortName+alltrim(str(x-1)) 
         hb_comSetDevice( nPort, cPort)
         IF hb_comOpen( nPort )
            aadd(pa,x)
            hb_comClose( nPort )
         ENDIF

      ENDIF
      x--
   enddo
return (pa)
