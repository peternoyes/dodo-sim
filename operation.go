package dodosim

import (
	//"fmt"
	"reflect"
	"runtime"
)

type Operation struct {
	Handler func(r Resolve) (bool, uint8)
	Mode    AddressMode
	Cycles  uint8
}

var table [256]*Operation

func BuildTable() {
	// brk,  ora,  nop,  slo,  nop,  ora,  asl,  slo,  php,  ora,  asl,  nop,  nop,  ora,  asl,
	// imp, indx,  imp, indx,   zp,   zp,   zp,   zp,  imp,  imm,  acc,  imm, abso, abso, abso, abso
	// 7,    6,    2,    8,    3,    3,    5,    5,    3,    2,    2,    2,    4,    4,    6,    6
	table[0] = &Operation{Brk, Imp, 7}
	table[1] = &Operation{Ora, Indx, 6}
	table[2] = &Operation{Nop, Imp, 2}
	table[3] = &Operation{Slo, Indx, 8}
	table[4] = &Operation{Nop, Zp, 3}
	table[5] = &Operation{Ora, Zp, 3}
	table[6] = &Operation{Asl, Zp, 5}
	table[7] = &Operation{Slo, Zp, 5}
	table[8] = &Operation{Php, Imp, 3}
	table[9] = &Operation{Ora, Imm, 2}
	table[10] = &Operation{Asl, Acc, 2}
	table[11] = &Operation{Nop, Imm, 2}
	table[12] = &Operation{Nop, Abso, 4}
	table[13] = &Operation{Ora, Abso, 4}
	table[14] = &Operation{Asl, Abso, 6}
	table[15] = &Operation{Slo, Abso, 6}

	// bpl,  ora,  nop,  slo,  nop,  ora,  asl,  slo,  clc,  ora,  nop,  slo,  nop,  ora,  asl,  slo
	// rel, indy,  imp, indy,  zpx,  zpx,  zpx,  zpx,  imp, absy,  imp, absy, absx, absx, absx, absx
	// 2,    5,    2,    8,    4,    4,    6,    6,    2,    4,    2,    7,    4,    4,    7,    7
	table[16] = &Operation{Bpl, Rel, 2}
	table[17] = &Operation{Ora, Indy, 5}
	table[18] = &Operation{Nop, Imp, 2}
	table[19] = &Operation{Slo, Indy, 8}
	table[20] = &Operation{Nop, Zpx, 4}
	table[21] = &Operation{Ora, Zpx, 4}
	table[22] = &Operation{Asl, Zpx, 6}
	table[23] = &Operation{Slo, Zpx, 6}
	table[24] = &Operation{Clc, Imp, 2}
	table[25] = &Operation{Ora, Absy, 4}
	table[26] = &Operation{Inc, Acc, 2}
	table[27] = &Operation{Slo, Absy, 7}
	table[28] = &Operation{Nop, Absx, 4}
	table[29] = &Operation{Ora, Absx, 4}
	table[30] = &Operation{Asl, Absx, 7}
	table[31] = &Operation{Slo, Absx, 7}

	// jsr,  and,  nop,  rla,  bit,  and,  rol,  rla,  plp,  and,  rol,  nop,  bit,  and,  rol,  rla
	// abso, indx,  imp, indx,   zp,   zp,   zp,   zp,  imp,  imm,  acc,  imm, abso, abso, abso, abso
	// 6,    6,    2,    8,    3,    3,    5,    5,    4,    2,    2,    2,    4,    4,    6,    6
	table[32] = &Operation{Jsr, Abso, 6}
	table[33] = &Operation{And, Indx, 6}
	table[34] = &Operation{Nop, Imp, 2}
	table[35] = &Operation{Rla, Indx, 8}
	table[36] = &Operation{Bit, Zp, 3}
	table[37] = &Operation{And, Zp, 3}
	table[38] = &Operation{Rol, Zp, 5}
	table[39] = &Operation{Rla, Zp, 5}
	table[40] = &Operation{Plp, Imp, 4}
	table[41] = &Operation{And, Imm, 2}
	table[42] = &Operation{Rol, Acc, 2}
	table[43] = &Operation{Nop, Imm, 2}
	table[44] = &Operation{Bit, Abso, 4}
	table[45] = &Operation{And, Abso, 4}
	table[46] = &Operation{Rol, Abso, 6}
	table[47] = &Operation{Rla, Abso, 6}

	// bmi,  and,  and,  rla,  nop,  and,  rol,  rla,  sec,  and,  nop,  rla,  nop,  and,  rol,  rla
	// rel, indy,  zp, indy,  zpx,  zpx,  zpx,  zpx,  imp, absy,  imp, absy, absx, absx, absx, absx
	// 2,    5,    5,    8,    4,    4,    6,    6,    2,    4,    2,    7,    4,    4,    7,    7
	table[48] = &Operation{Bmi, Rel, 2}
	table[49] = &Operation{And, Indy, 5}
	table[50] = &Operation{And, Indzp, 5}
	table[51] = &Operation{Rla, Indy, 8}
	table[52] = &Operation{Nop, Zpx, 4}
	table[53] = &Operation{And, Zpx, 4}
	table[54] = &Operation{Rol, Zpx, 6}
	table[55] = &Operation{Rla, Zpx, 6}
	table[56] = &Operation{Sec, Imp, 2}
	table[57] = &Operation{And, Absy, 4}
	table[58] = &Operation{Nop, Imp, 2}
	table[59] = &Operation{Rla, Absy, 7}
	table[60] = &Operation{Nop, Absx, 4}
	table[61] = &Operation{And, Absx, 4}
	table[62] = &Operation{Rol, Absx, 7}
	table[63] = &Operation{Rla, Absx, 7}

	// rti,  eor,  nop,  sre,  nop,  eor,  lsr,  sre,  pha,  eor,  lsr,  nop,  jmp,  eor,  lsr, sre
	// imp, indx,  imp, indx,   zp,   zp,   zp,   zp,  imp,  imm,  acc,  imm, abso, abso, abso, abso
	// 6,    6,    2,    8,    3,    3,    5,    5,    3,    2,    2,    2,    3,    4,    6,    6
	table[64] = &Operation{Rti, Imp, 6}
	table[65] = &Operation{Eor, Indx, 6}
	table[66] = &Operation{Nop, Imp, 2}
	table[67] = &Operation{Sre, Indx, 8}
	table[68] = &Operation{Nop, Zp, 3}
	table[69] = &Operation{Eor, Zp, 3}
	table[70] = &Operation{Lsr, Zp, 5}
	table[71] = &Operation{Sre, Zp, 5}
	table[72] = &Operation{Pha, Imp, 3}
	table[73] = &Operation{Eor, Imm, 2}
	table[74] = &Operation{Lsr, Acc, 2}
	table[75] = &Operation{Nop, Imm, 2}
	table[76] = &Operation{Jmp, Abso, 3}
	table[77] = &Operation{Eor, Abso, 4}
	table[78] = &Operation{Lsr, Abso, 6}
	table[79] = &Operation{Sre, Abso, 6}

	// bvc,  eor,  nop,  sre,  nop,  eor,  lsr,  sre,  cli,  eor,  nop,  sre,  nop,  eor,  lsr,  sre
	// rel, indy,  imp, indy,  zpx,  zpx,  zpx,  zpx,  imp, absy,  imp, absy, absx, absx, absx, absx
	// 2,    5,    2,    8,    4,    4,    6,    6,    2,    4,    2,    7,    4,    4,    7,    7
	table[80] = &Operation{Bvc, Rel, 2}
	table[81] = &Operation{Eor, Indy, 5}
	table[82] = &Operation{Nop, Imp, 2}
	table[83] = &Operation{Sre, Indy, 8}
	table[84] = &Operation{Nop, Zpx, 4}
	table[85] = &Operation{Eor, Zpx, 4}
	table[86] = &Operation{Lsr, Zpx, 6}
	table[87] = &Operation{Sre, Zpx, 6}
	table[88] = &Operation{Cli, Imp, 2}
	table[89] = &Operation{Eor, Absy, 4}
	table[90] = &Operation{Nop, Imp, 2}
	table[91] = &Operation{Sre, Absy, 7}
	table[92] = &Operation{Nop, Absx, 4}
	table[93] = &Operation{Eor, Absx, 4}
	table[94] = &Operation{Lsr, Absx, 7}
	table[95] = &Operation{Sre, Absx, 7}

	// rts,  adc,  nop,  rra,  nop,  adc,  ror,  rra,  pla,  adc,  ror,  nop,  jmp,  adc,  ror, rra
	// imp, indx,  imp, indx,   zp,   zp,   zp,   zp,  imp,  imm,  acc,  imm,  ind, abso, abso, abso
	// 6,    6,    2,    8,    3,    3,    5,    5,    4,    2,    2,    2,    5,    4,    6,    6
	table[96] = &Operation{Rts, Imp, 6}
	table[97] = &Operation{Adc, Indx, 6}
	table[98] = &Operation{Nop, Imp, 2}
	table[99] = &Operation{Rra, Indx, 8}
	table[100] = &Operation{Nop, Zp, 3}
	table[101] = &Operation{Adc, Zp, 3}
	table[102] = &Operation{Ror, Zp, 5}
	table[103] = &Operation{Rra, Zp, 5}
	table[104] = &Operation{Pla, Imp, 4}
	table[105] = &Operation{Adc, Imm, 2}
	table[106] = &Operation{Ror, Acc, 2}
	table[107] = &Operation{Nop, Imm, 2}
	table[108] = &Operation{Jmp, Ind, 5}
	table[109] = &Operation{Adc, Abso, 4}
	table[110] = &Operation{Ror, Abso, 6}
	table[111] = &Operation{Rra, Abso, 6}

	// bvs,  adc,  nop,  rra,  nop,  adc,  ror,  rra,  sei,  adc,  adc,  rra,  nop,  adc,  ror, rra
	// rel, indy,  imp, indy,  zpx,  zpx,  zpx,  zpx,  imp, absy,  zp, absy, absx, absx, absx, absx
	// 2,    5,    2,    8,    4,    4,    6,    6,    2,    4,    2,    7,    4,    4,    7,    7
	table[112] = &Operation{Bvs, Rel, 2}
	table[113] = &Operation{Adc, Indy, 5}
	table[114] = &Operation{Adc, Indzp, 5}
	table[115] = &Operation{Rra, Indy, 8}
	table[116] = &Operation{Nop, Zpx, 4}
	table[117] = &Operation{Adc, Zpx, 4}
	table[118] = &Operation{Ror, Zpx, 6}
	table[119] = &Operation{Rra, Zpx, 6}
	table[120] = &Operation{Sei, Imp, 2}
	table[121] = &Operation{Adc, Absy, 4}
	table[122] = &Operation{Nop, Imp, 2}
	table[123] = &Operation{Rra, Absy, 7}
	table[124] = &Operation{Nop, Absx, 4}
	table[125] = &Operation{Adc, Absx, 4}
	table[126] = &Operation{Ror, Absx, 7}
	table[127] = &Operation{Rra, Absx, 7}

	// nop,  sta,  nop,  sax,  sty,  sta,  stx,  sax,  dey,  nop,  txa,  nop,  sty,  sta,  stx,  sax
	// imm, indx,  imm, indx,   zp,   zp,   zp,   zp,  imp,  imm,  imp,  imm, abso, abso, abso,  abso
	// 2,    6,    2,    6,    3,    3,    3,    3,    2,    2,    2,    2,    4,    4,    4,    4
	table[128] = &Operation{Nop, Imm, 2}
	table[129] = &Operation{Sta, Indx, 6}
	table[130] = &Operation{Nop, Imm, 2}
	table[131] = &Operation{Sax, Indx, 6}
	table[132] = &Operation{Sty, Zp, 3}
	table[133] = &Operation{Sta, Zp, 3}
	table[134] = &Operation{Stx, Zp, 3}
	table[135] = &Operation{Sax, Zp, 3}
	table[136] = &Operation{Dey, Imp, 2}
	table[137] = &Operation{Nop, Imm, 2}
	table[138] = &Operation{Txa, Imp, 2}
	table[139] = &Operation{Nop, Imm, 2}
	table[140] = &Operation{Sty, Abso, 4}
	table[141] = &Operation{Sta, Abso, 4}
	table[142] = &Operation{Stx, Abso, 4}
	table[143] = &Operation{Sax, Abso, 4}

	// bcc,  sta,  nop,  nop,  sty,  sta,  stx,  sax,  tya,  sta,  txs,  nop,  nop,  sta,  nop,  nop
	// rel, indy,  imp, indy,  zpx,  zpx,  zpy,  zpy,  imp, absy,  imp, absy, absx, absx, absy, absy
	// 2,    6,    2,    6,    4,    4,    4,    4,    2,    5,    2,    5,    5,    5,    5,    5
	table[144] = &Operation{Bcc, Rel, 2}
	table[145] = &Operation{Sta, Indy, 6}
	table[146] = &Operation{Nop, Imp, 2}
	table[147] = &Operation{Nop, Indy, 6}
	table[148] = &Operation{Sty, Zpx, 4}
	table[149] = &Operation{Sta, Zpx, 4}
	table[150] = &Operation{Stx, Zpy, 4}
	table[151] = &Operation{Sax, Zpy, 4}
	table[152] = &Operation{Tya, Imp, 2}
	table[153] = &Operation{Sta, Absy, 5}
	table[154] = &Operation{Txs, Imp, 2}
	table[155] = &Operation{Nop, Absy, 5}
	table[156] = &Operation{Nop, Absx, 5}
	table[157] = &Operation{Sta, Absx, 5}
	table[158] = &Operation{Nop, Absy, 5}
	table[159] = &Operation{Nop, Absy, 5}

	// ldy,  lda,  ldx,  lax,  ldy,  lda,  ldx,  lax,  tay,  lda,  tax,  nop,  ldy,  lda,  ldx,  lax
	// imm, indx,  imm, indx,   zp,   zp,   zp,   zp,  imp,  imm,  imp,  imm, abso, abso, abso,  abso
	// 2,    6,    2,    6,    3,    3,    3,    3,    2,    2,    2,    2,    4,    4,    4,    4
	table[160] = &Operation{Ldy, Imm, 2}
	table[161] = &Operation{Lda, Indx, 6}
	table[162] = &Operation{Ldx, Imm, 2}
	table[163] = &Operation{Lax, Indx, 6}
	table[164] = &Operation{Ldy, Zp, 3}
	table[165] = &Operation{Lda, Zp, 3}
	table[166] = &Operation{Ldx, Zp, 3}
	table[167] = &Operation{Lax, Zp, 3}
	table[168] = &Operation{Tay, Imp, 2}
	table[169] = &Operation{Lda, Imm, 2}
	table[170] = &Operation{Tax, Imp, 2}
	table[171] = &Operation{Nop, Imm, 2}
	table[172] = &Operation{Ldy, Abso, 4}
	table[173] = &Operation{Lda, Abso, 4}
	table[174] = &Operation{Ldx, Abso, 4}
	table[175] = &Operation{Lax, Abso, 4}

	// bcs,  lda,  lda,  lax,  ldy,  lda,  ldx,  lax,  clv,  lda,  tsx,  lax,  ldy,  lda,  ldx, lax
	// rel, indy,  zp, indy,  zpx,  zpx,  zpy,  zpy,  imp, absy,  imp, absy, absx, absx, absy, absy
	// 2,    5,    5,    5,    4,    4,    4,    4,    2,    4,    2,    4,    4,    4,    4,    4
	table[176] = &Operation{Bcs, Rel, 2}
	table[177] = &Operation{Lda, Indy, 5}
	table[178] = &Operation{Lda, Indzp, 5}
	table[179] = &Operation{Lax, Indy, 5}
	table[180] = &Operation{Ldy, Zpx, 4}
	table[181] = &Operation{Lda, Zpx, 4}
	table[182] = &Operation{Ldx, Zpy, 4}
	table[183] = &Operation{Lax, Zpy, 4}
	table[184] = &Operation{Clv, Imp, 2}
	table[185] = &Operation{Lda, Absy, 4}
	table[186] = &Operation{Tsx, Imp, 2}
	table[187] = &Operation{Lax, Absy, 4}
	table[188] = &Operation{Ldy, Absx, 4}
	table[189] = &Operation{Lda, Absx, 4}
	table[190] = &Operation{Ldx, Absy, 4}
	table[191] = &Operation{Lax, Absy, 4}

	// cpy,  cmp,  nop,  dcp,  cpy,  cmp,  dec,  dcp,  iny,  cmp,  dex,  nop,  cpy,  cmp,  dec,  dcp
	// imm, indx,  imm, indx,   zp,   zp,   zp,   zp,  imp,  imm,  imp,  imm, abso, abso, abso,  abso
	// 2,    6,    2,    8,    3,    3,    5,    5,    2,    2,    2,    2,    4,    4,    6,    6
	table[192] = &Operation{Cpy, Imm, 2}
	table[193] = &Operation{Cmp, Indx, 6}
	table[194] = &Operation{Nop, Imm, 2}
	table[195] = &Operation{Dcp, Indx, 8}
	table[196] = &Operation{Cpy, Zp, 3}
	table[197] = &Operation{Cmp, Zp, 3}
	table[198] = &Operation{Dec, Zp, 5}
	table[199] = &Operation{Dcp, Zp, 5}
	table[200] = &Operation{Iny, Imp, 2}
	table[201] = &Operation{Cmp, Imm, 2}
	table[202] = &Operation{Dex, Imp, 2}
	table[203] = &Operation{Nop, Imm, 2}
	table[204] = &Operation{Cpy, Abso, 4}
	table[205] = &Operation{Cmp, Abso, 4}
	table[206] = &Operation{Dec, Abso, 6}
	table[207] = &Operation{Dcp, Abso, 6}

	// bne,  cmp,  nop,  dcp,  nop,  cmp,  dec,  dcp,  cld,  cmp,  nop,  dcp,  nop,  cmp,  dec,  dcp
	// rel, indy,  imp, indy,  zpx,  zpx,  zpx,  zpx,  imp, absy,  imp, absy, absx, absx, absx,  absx
	// 2,    5,    2,    8,    4,    4,    6,    6,    2,    4,    2,    7,    4,    4,    7,    7
	table[208] = &Operation{Bne, Rel, 2}
	table[209] = &Operation{Cmp, Indy, 5}
	table[210] = &Operation{Nop, Imp, 2}
	table[211] = &Operation{Dcp, Indy, 8}
	table[212] = &Operation{Nop, Zpx, 4}
	table[213] = &Operation{Cmp, Zpx, 4}
	table[214] = &Operation{Dec, Zpx, 6}
	table[215] = &Operation{Dcp, Zpx, 6}
	table[216] = &Operation{Cld, Imp, 2}
	table[217] = &Operation{Cmp, Absy, 4}
	table[218] = &Operation{Nop, Imp, 2}
	table[219] = &Operation{Dcp, Absy, 7}
	table[220] = &Operation{Nop, Absx, 4}
	table[221] = &Operation{Cmp, Absx, 4}
	table[222] = &Operation{Dec, Absx, 7}
	table[223] = &Operation{Dcp, Absx, 7}

	// cpx,  sbc,  nop,  isb,  cpx,  sbc,  inc,  isb,  inx,  sbc,  nop,  sbc,  cpx,  sbc,  inc,  isb
	// imm, indx,  imm, indx,   zp,   zp,   zp,   zp,  imp,  imm,  imp,  imm, abso, abso, abso, abso
	// 2,    6,    2,    8,    3,    3,    5,    5,    2,    2,    2,    2,    4,    4,    6,    6
	table[224] = &Operation{Cpx, Imm, 2}
	table[225] = &Operation{Sbc, Indx, 6}
	table[226] = &Operation{Nop, Imm, 2}
	table[227] = &Operation{Isb, Indx, 8}
	table[228] = &Operation{Cpx, Zp, 3}
	table[229] = &Operation{Sbc, Zp, 3}
	table[230] = &Operation{Inc, Zp, 5}
	table[231] = &Operation{Isb, Zp, 5}
	table[232] = &Operation{Inx, Imp, 2}
	table[233] = &Operation{Sbc, Imm, 2}
	table[234] = &Operation{Nop, Imp, 2}
	table[235] = &Operation{Sbc, Imm, 2}
	table[236] = &Operation{Cpx, Abso, 4}
	table[237] = &Operation{Sbc, Abso, 4}
	table[238] = &Operation{Inc, Abso, 6}
	table[239] = &Operation{Isb, Abso, 6}

	// beq,  sbc,  nop,  isb,  nop,  sbc,  inc,  isb,  sed,  sbc,  nop,  isb,  nop,  sbc,  inc,  isb
	// rel, indy,  imp, indy,  zpx,  zpx,  zpx,  zpx,  imp, absy,  imp, absy, absx, absx, absx, absx
	// 2,    5,    2,    8,    4,    4,    6,    6,    2,    4,    2,    7,    4,    4,    7,    7
	table[240] = &Operation{Beq, Rel, 4}
	table[241] = &Operation{Sbc, Indy, 5}
	table[242] = &Operation{Nop, Imp, 2}
	table[243] = &Operation{Isb, Indy, 8}
	table[244] = &Operation{Nop, Zpx, 4}
	table[245] = &Operation{Sbc, Zpx, 4}
	table[246] = &Operation{Inc, Zpx, 6}
	table[247] = &Operation{Isb, Zpx, 6}
	table[248] = &Operation{Sed, Imp, 2}
	table[249] = &Operation{Sbc, Absy, 4}
	table[250] = &Operation{Nop, Imp, 2}
	table[251] = &Operation{Isb, Absy, 7}
	table[252] = &Operation{Nop, Absx, 4}
	table[253] = &Operation{Sbc, Absx, 4}
	table[254] = &Operation{Inc, Absx, 7}
	table[255] = &Operation{Isb, Absx, 7}
}

func GetOperation(opcode uint8) *Operation {
	return table[opcode]
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (o *Operation) Execute(cpu *Cpu, space Space, opcode uint8) uint8 {
	c := o.Cycles

	//fmt.Println(GetFunctionName(o.Handler), " Mode: ", o.Mode)

	r := o.Mode.Resolve(cpu, space, opcode)
	pop, rc := o.Handler(r)
	c += rc
	if r.Penalty && pop {
		c += 1
	}

	return c
}

func Adc(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var c uint8 = 0
	var res uint16

	res = uint16(cpu.A) + val + uint16(cpu.Status&Carry)
	cpu.ZeroCalc(res)

	if cpu.Status&Decimal != 0 {

		if ((cpu.A & 0xF) + (uint8(val) & 0xF) + uint8(cpu.Status&Carry)) > 0x9 {
			res += 0x6
		}

		cpu.SignCalc(res)
		cpu.OverflowCalc(res, cpu.A, val)

		if res > 0x99 {
			res += 96
		}

		if res > 0x99 {
			cpu.SetCarry()
		} else {
			cpu.ClearCarry()
		}

		c = 1 // One extra clock cycle for BDC
	} else {
		cpu.CarryCalc(res)
		cpu.OverflowCalc(res, cpu.A, val)
		cpu.SignCalc(res)
	}

	cpu.SaveAccum(res)

	return true, c // Possible penalty
}

func And(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = uint16(cpu.A) & val

	cpu.ZeroCalc(res)
	cpu.SignCalc(res)

	cpu.SaveAccum(res)

	return true, 0
}

func Asl(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var res uint16
	res = r.Read() << 1

	cpu.CarryCalc(res)
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)

	r.Write(res)

	return false, 0 // No potential for penalty
}

func Bcc(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Carry == 0 {
		oldpc := cpu.PC
		cpu.PC += r.Address

		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Bcs(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Carry == Carry {
		oldpc := cpu.PC
		cpu.PC += r.Address

		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Beq(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Zero == Zero {
		oldpc := cpu.PC
		cpu.PC += r.Address
		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Bit(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = uint16(cpu.A) & val

	cpu.ZeroCalc(res)
	cpu.Status = Status((uint8(cpu.Status) & 0x3F) | uint8(val&0xC0))
	return false, 0
}

func Bmi(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Sign == Sign {
		oldpc := cpu.PC
		cpu.PC += r.Address
		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Bne(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Zero == 0 {
		oldpc := cpu.PC
		cpu.PC += r.Address
		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Bpl(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Sign == 0 {
		oldpc := cpu.PC
		cpu.PC += r.Address
		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Brk(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.PC++
	r.Push16(cpu.PC)
	r.Push8(uint8(cpu.Status | Break))
	cpu.SetInterrupt()
	cpu.PC = uint16(r.Space.Read(0xFFFE)) | (uint16(r.Space.Read(0xFFFF)) << 8)
	return false, 0
}

func Bvc(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Overflow == 0 {
		oldpc := cpu.PC
		cpu.PC += r.Address
		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Bvs(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var c uint8 = 0
	if cpu.Status&Overflow == Overflow {
		oldpc := cpu.PC
		cpu.PC += r.Address
		if (oldpc & 0xFF00) != (cpu.PC & 0xFF00) {
			c = 2
		} else {
			c = 1
		}
	}
	return false, c
}

func Clc(r Resolve) (bool, uint8) {
	r.Cpu.ClearCarry()
	return false, 0
}

func Cld(r Resolve) (bool, uint8) {
	r.Cpu.ClearDecimal()
	return false, 0
}

func Cli(r Resolve) (bool, uint8) {
	r.Cpu.ClearInterrupt()
	return false, 0
}

func Clv(r Resolve) (bool, uint8) {
	r.Cpu.ClearOverflow()
	return false, 0
}

func Cmp(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()

	//fmt.Println("Cmp: ", val)

	var res uint16
	res = uint16(cpu.A) - val

	if cpu.A >= uint8(val&0x00FF) {
		cpu.SetCarry()
	} else {
		cpu.ClearCarry()
	}

	if cpu.A == uint8(val&0x00FF) {
		cpu.SetZero()
	} else {
		cpu.ClearZero()
	}

	cpu.SignCalc(res)

	return true, 0
}

func Cpx(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = uint16(cpu.X) - val

	if cpu.X >= uint8(val&0x00FF) {
		cpu.SetCarry()
	} else {
		cpu.ClearCarry()
	}

	if cpu.X == uint8(val&0x00FF) {
		cpu.SetZero()
	} else {
		cpu.ClearZero()
	}

	cpu.SignCalc(res)

	return false, 0
}

func Cpy(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = uint16(cpu.Y) - val

	if cpu.Y >= uint8(val&0x00FF) {
		cpu.SetCarry()
	} else {
		cpu.ClearCarry()
	}

	if cpu.Y == uint8(val&0x00FF) {
		cpu.SetZero()
	} else {
		cpu.ClearZero()
	}

	cpu.SignCalc(res)

	return false, 0
}

func Dec(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var res uint16
	res = r.Read() - 1
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)
	r.Write(res)
	return false, 0
}

func Dex(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.X -= 1
	cpu.ZeroCalc8(cpu.X)
	cpu.SignCalc8(cpu.X)
	return false, 0
}

func Dey(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.Y -= 1
	cpu.ZeroCalc8(cpu.Y)
	cpu.SignCalc8(cpu.Y)
	return false, 0
}

func Eor(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var res uint16
	res = uint16(cpu.A) ^ r.Read()
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)
	cpu.SaveAccum(res)
	return true, 0
}

func Inc(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var res uint16
	res = r.Read() + 1
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)
	r.Write(res)
	return false, 0
}

func Inx(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.X++
	cpu.ZeroCalc(uint16(cpu.X))
	cpu.SignCalc(uint16(cpu.X))
	return false, 0
}

func Iny(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.Y++
	cpu.ZeroCalc(uint16(cpu.Y))
	cpu.SignCalc(uint16(cpu.Y))
	return false, 0
}

func Jmp(r Resolve) (bool, uint8) {
	r.Cpu.PC = r.Address
	return false, 0
}

func Jsr(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	r.Push16(cpu.PC - 1)
	cpu.PC = r.Address
	return false, 0
}

func Lda(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var res uint16
	res = r.Read() & 0x00FF
	cpu.A = uint8(res)
	cpu.ZeroCalc8(cpu.A)
	cpu.SignCalc8(cpu.A)
	return true, 0
}

func Ldx(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var res uint16
	res = r.Read() & 0x00FF
	cpu.X = uint8(res)
	cpu.ZeroCalc8(cpu.X)
	cpu.SignCalc8(cpu.X)
	return true, 0
}

func Ldy(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	var res uint16
	res = r.Read() & 0x00FF
	cpu.Y = uint8(res)
	cpu.ZeroCalc8(cpu.Y)
	cpu.SignCalc8(cpu.Y)
	return true, 0
}

func Lsr(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = val >> 1

	if val&1 != 0 {
		cpu.SetCarry()
	} else {
		cpu.ClearCarry()
	}
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)

	r.Write(res)

	return false, 0
}

func Nop(r Resolve) (bool, uint8) {
	//fmt.Println("NOP: ", r.Opcode)

	switch r.Opcode {
	default:
		return false, 0
	case 0x1C:
		fallthrough
	case 0x3C:
		fallthrough
	case 0x5C:
		fallthrough
	case 0x7C:
		fallthrough
	case 0xDC:
		fallthrough
	case 0xFC:
		return true, 0
	}
}

func Ora(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = uint16(cpu.A) | val
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)
	cpu.SaveAccum(res)
	return true, 0
}

func Pha(r Resolve) (bool, uint8) {
	r.Push8(r.Cpu.A)
	return false, 0
}

func Php(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	r.Push8(uint8(cpu.Status | Break))
	return false, 0
}

func Pla(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.A = r.Pull8()
	cpu.ZeroCalc(uint16(cpu.A))
	cpu.SignCalc(uint16(cpu.A))
	return false, 0
}

func Plp(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.Status = Status(r.Pull8()) | Constant
	return false, 0
}

func Rol(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = (val << 1) | uint16(cpu.Status&Carry)

	cpu.CarryCalc(res)
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)

	r.Write(res)

	return false, 0
}

func Ror(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	var res uint16
	res = (val >> 1) | uint16(uint8(cpu.Status&Carry)<<7)
	if val&1 != 0 {
		cpu.SetCarry()
	} else {
		cpu.ClearCarry()
	}
	cpu.ZeroCalc(res)
	cpu.SignCalc(res)

	r.Write(res)

	return false, 0
}

func Rti(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.Status = Status(r.Pull8())
	cpu.PC = r.Pull16()
	return false, 0
}

func Rts(r Resolve) (bool, uint8) {
	val := r.Pull16()
	r.Cpu.PC = val + 1
	return false, 0
}

func Sbc(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	val := r.Read()
	carry := uint8(1)
	if cpu.Status&Carry != 0 {
		carry = 0
	}

	var res uint16
	var c uint8 = 0
	if cpu.Status&Decimal != 0 {
		low := (cpu.A & 0x0F) - (uint8(val) & 0x0F) - carry
		if (low & 0x10) != 0 {
			low -= 6
		}
		nextCarry := uint8(0)
		if (low & 0x10) != 0 {
			nextCarry = 1
		}

		high := (cpu.A >> 4) - (uint8(val) >> 4) - nextCarry
		if (high & 0x10) != 0 {
			high -= 6
		}

		res = uint16(low&0x0F) | uint16((high<<4)&0xF0)

		if (high & 0xFF) < 15 {
			cpu.SetCarry()
		} else {
			cpu.ClearCarry()
		}
		if res == 0 {
			cpu.SetZero()
		} else {
			cpu.ClearZero()
		}
		cpu.ClearSign()
		cpu.ClearOverflow()

		c = 1 // One extra clock cycle for BDC
	} else {

		res = uint16(cpu.A) - val - uint16(carry)

		cpu.SignCalc(res)
		cpu.ZeroCalc(res)

		if ((uint16(cpu.A)^res)&0x80) != 0 && ((uint16(cpu.A)^val)&0x80) != 0 {
			cpu.SetOverflow()
		} else {
			cpu.ClearOverflow()
		}

		if res < 0x100 {
			cpu.SetCarry()
		} else {
			cpu.ClearCarry()
		}
	}

	cpu.SaveAccum(res)

	return true, c // Possible penalty
}

func Sec(r Resolve) (bool, uint8) {
	r.Cpu.SetCarry()
	return false, 0
}

func Sed(r Resolve) (bool, uint8) {
	r.Cpu.SetDecimal()
	return false, 0
}

func Sei(r Resolve) (bool, uint8) {
	r.Cpu.SetInterrupt()
	return false, 0
}

func Sta(r Resolve) (bool, uint8) {
	r.Write(uint16(r.Cpu.A))
	return false, 0
}

func Stx(r Resolve) (bool, uint8) {
	r.Write(uint16(r.Cpu.X))
	return false, 0
}

func Sty(r Resolve) (bool, uint8) {
	r.Write(uint16(r.Cpu.Y))
	return false, 0
}

func Tax(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.X = cpu.A
	cpu.ZeroCalc(uint16(cpu.X))
	cpu.SignCalc(uint16(cpu.X))
	return false, 0
}

func Tay(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.Y = cpu.A
	cpu.ZeroCalc(uint16(cpu.Y))
	cpu.SignCalc(uint16(cpu.Y))
	return false, 0
}

func Tsx(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.X = cpu.SP
	cpu.ZeroCalc(uint16(cpu.X))
	cpu.SignCalc(uint16(cpu.X))
	return false, 0
}

func Txa(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.A = cpu.X
	cpu.ZeroCalc(uint16(cpu.A))
	cpu.SignCalc(uint16(cpu.A))
	return false, 0
}

func Txs(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.SP = cpu.X
	return false, 0
}

func Tya(r Resolve) (bool, uint8) {
	cpu := r.Cpu
	cpu.A = cpu.Y
	cpu.ZeroCalc(uint16(cpu.A))
	cpu.SignCalc(uint16(cpu.A))
	return false, 0
}

func Lax(r Resolve) (bool, uint8) {
	panic("Lax Not Implemented")
}

func Sax(r Resolve) (bool, uint8) {
	panic("Sax Not Implemented")
}

func Dcp(r Resolve) (bool, uint8) {
	panic("Dcp Not Implemented")
}

func Isb(r Resolve) (bool, uint8) {
	panic("Isb Not Implemented")
}

func Slo(r Resolve) (bool, uint8) {
	panic("Slo Not Implemented")
}

func Rla(r Resolve) (bool, uint8) {
	panic("Rla Not Implemented")
}

func Sre(r Resolve) (bool, uint8) {
	panic("Sre Not Implemented")
}

func Rra(r Resolve) (bool, uint8) {
	panic("Rra Not Implemented")
}
