// +build avr,atmega

package runtime

import (
	"device/avr"
	"machine"
)

func initUART() {
	machine.UART0.Configure(machine.UARTConfig{})
}

func putchar(c byte) {
	machine.UART0.WriteByte(c)
}

// Sleep for a given period. The period is defined by the WDT peripheral, and is
// on most chips (at least) 3 bits wide, in powers of two from 16ms to 2s
// (0=16ms, 1=32ms, 2=64ms...). Note that the WDT is not very accurate: it can
// be off by a large margin depending on temperature and supply voltage.
//
// TODO: disable more peripherals etc. to reduce sleep current.
func sleepWDT(period uint8) {
	// Configure WDT
	avr.Asm("cli")
	avr.Asm("wdr")
	// Start timed sequence.
	avr.WDTCSR.SetBits(avr.WDTCSR_WDCE | avr.WDTCSR_WDE)
	// Enable WDT and set new timeout
	avr.WDTCSR.SetBits(avr.WDTCSR_WDIE | period)
	avr.Asm("sei")

	// Set sleep mode to idle and enable sleep mode.
	// Note: when using something other than idle, the UART won't work
	// correctly. This needs to be fixed, though, so we can truly sleep.
	avr.SMCR.Set((0 << 1) | avr.SMCR_SE)

	// go to sleep
	avr.Asm("sleep")

	// disable sleep
	avr.SMCR.Set(0)
}
