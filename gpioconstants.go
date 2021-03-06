package golibgpiod

const GPIOD_LINE_DIRECTION_INPUT int = 1
const GPIOD_LINE_DIRECTION_OUTPUT int = 2

const GPIOD_LINE_ACTIVE_STATE_HIGH int = 1
const GPIOD_LINE_ACTIVE_STATE_LOW int = 2

const GPIOLINE_FLAG_KERNEL uint32 = 1
const GPIOLINE_FLAG_IS_OUT uint32 = 1 << 1
const GPIOLINE_FLAG_ACTIVE_LOW uint32 = 1 << 2
const GPIOLINE_FLAG_OPEN_DRAIN uint32 = 1 << 3
const GPIOLINE_FLAG_OPEN_SOURCE uint32 = 1 << 4
