        text    "".Str2Bytes1(SB), NOSPLIT|ABIInternal, $0-40
        pcdata  $0, $-2
        pcdata  $1, $-2
        funcdata        $0, gclocals·39d1b96ca581879f548ad2c8aeb3a5fe(SB)
        funcdata        $1, gclocals·7d2d5fca80364273fb07d5820a76fef4(SB)
        funcdata        $2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        funcdata        $3, "".Str2Bytes1.stkobj(SB)
        pcdata  $0, $0
        pcdata  $1, $1
        
        // zero 'b'
        movq    $0, "".b+24(SP)
        xorps   X0, X0
        movups  X0, "".b+32(SP)
        
        // bh.Data = sh.Data
        movq    "".s+8(SP), AX
        movq    AX, "".b+24(SP)
        // bh.Cap = sh.Len
        movq    "".s+16(SP), AX
        movq    AX, "".b+40(SP)
        
        // bh.Len = sh.Len
        pcdata  $1, $2
        movq    "".s+16(SP), AX
        movq    AX, "".b+32(SP)
        ret


        text    "".Str2Bytes2(SB), NOSPLIT|ABIInternal, $0-40
        pcdata  $0, $-2
        pcdata  $1, $-2
        funcdata        $0, gclocals·39d1b96ca581879f548ad2c8aeb3a5fe(SB)
        funcdata        $1, gclocals·7d2d5fca80364273fb07d5820a76fef4(SB)
        funcdata        $2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        funcdata        $3, "".Str2Bytes2.stkobj(SB)
        pcdata  $0, $0
        pcdata  $1, $1
        movq    $0, "".b+24(SP)
        xorps   X0, X0
        movups  X0, "".b+32(SP)
        movq    "".s+8(SP), AX
        movq    AX, "".b+24(SP)
        movq    "".s+16(SP), AX
        movq    AX, "".b+40(SP)
        pcdata  $1, $2
        movq    "".s+16(SP), AX
        movq    AX, "".b+32(SP)
        ret
        
        text    "".Str2Bytes3(SB), NOSPLIT|ABIInternal, $0-40
        pcdata  $0, $-2
        pcdata  $1, $-2
        funcdata        $0, gclocals·39d1b96ca581879f548ad2c8aeb3a5fe(SB)
        funcdata        $1, gclocals·7d2d5fca80364273fb07d5820a76fef4(SB)
        funcdata        $2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        funcdata        $3, "".Str2Bytes3.stkobj(SB)
        pcdata  $0, $0
        pcdata  $1, $1
        movq    $0, "".b+24(SP)
        xorps   X0, X0
        movups  X0, "".b+32(SP)
        movq    "".s+8(SP), AX
        movq    AX, "".b+24(SP)
        movq    "".s+16(SP), AX
        movq    AX, "".b+40(SP)
        pcdata  $1, $2
        movq    "".s+16(SP), AX
        movq    AX, "".b+32(SP)
        ret
        
        text    "".Str2Bytes4(SB), NOSPLIT|ABIInternal, $0-40
        pcdata  $0, $-2
        pcdata  $1, $-2
        funcdata        $0, gclocals·24c1b9f537195769589945d5dec53fc6(SB)
        funcdata        $1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
        funcdata        $2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        funcdata        $3, "".Str2Bytes4.stkobj(SB)
        pcdata  $0, $0
        pcdata  $1, $1
        // zero 'b'
        movq    $0, "".b+24(SP)
        xorps   X0, X0
        movups  X0, "".b+32(SP)
        // bh.Data = sh.Data
        movq    "".s+8(SP), AX
        movq    AX, "".b+24(SP)
        // bh.Cap = sh.Cap
        movq    "".s+16(SP), AX
        movq    AX, "".b+40(SP)
        // bh.Len = sh.Len
        movq    "".s+16(SP), AX
        movq    AX, "".b+32(SP)
        ret


        text    "".Str2Bytes5(SB), NOSPLIT|ABIInternal, $0-40
        pcdata  $0, $-2
        pcdata  $1, $-2
        funcdata        $0, gclocals·39d1b96ca581879f548ad2c8aeb3a5fe(SB)
        funcdata        $1, gclocals·7d2d5fca80364273fb07d5820a76fef4(SB)
        funcdata        $2, gclocals·9fb7f0986f647f17cb53dda1484e0f7a(SB)
        funcdata        $3, "".Str2Bytes5.stkobj(SB)
        pcdata  $0, $0
        pcdata  $1, $1
        // zero 'b'
        movq    $0, "".b+24(SP)
        xorps   X0, X0
        movups  X0, "".b+32(SP)
        // bh.Data = sh.Data
        pcdata  $0, $1
        movq    "".s+8(SP), AX
        pcdata  $0, $0
        movq    AX, "".b+24(SP)
        // bh.Cap = sh.Cap
        movq    "".s+16(SP), AX
        movq    AX, "".b+40(SP)
        // bh.Len = sh.Len
        pcdata  $1, $2
        movq    "".s+16(SP), AX
        movq    AX, "".b+32(SP)
        ret