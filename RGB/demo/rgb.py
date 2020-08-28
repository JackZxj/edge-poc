#!/usr/bin/env python3
# encoding: utf-8

import RPi.GPIO as gpio
import time

R, G, B = 12, 13, 19

gpio.setmode(gpio.BCM)

gpio.setup(R, gpio.OUT)
gpio.setup(G, gpio.OUT)
gpio.setup(B, gpio.OUT)

pwmR = gpio.PWM(R, 50)
pwmG = gpio.PWM(G, 50)
pwmB = gpio.PWM(B, 50)

pwmR.start(0)
pwmG.start(0)
pwmB.start(0)

try:
    t = 0.3
    while True:
        # 红色灯全亮，蓝灯，绿灯全暗（红色）
        # lights red
        pwmR.ChangeDutyCycle(100)
        pwmG.ChangeDutyCycle(0)
        pwmB.ChangeDutyCycle(0)
        time.sleep(t)

        # 绿色灯全亮，红灯，蓝灯全暗（绿色）
        # lights green
        pwmR.ChangeDutyCycle(0)
        pwmG.ChangeDutyCycle(100)
        pwmB.ChangeDutyCycle(0)
        time.sleep(t)

        # 蓝色灯全亮，红灯，绿灯全暗（蓝色）
        # lights blue
        pwmR.ChangeDutyCycle(0)
        pwmG.ChangeDutyCycle(0)
        pwmB.ChangeDutyCycle(100)
        time.sleep(t)

        # 红灯，绿灯全亮，蓝灯全暗（黄色）
        # lights red and green => yellow
        pwmR.ChangeDutyCycle(100)
        pwmG.ChangeDutyCycle(100)
        pwmB.ChangeDutyCycle(0)
        time.sleep(t)

        # 红灯，蓝灯全亮，绿灯全暗（洋红色）
        # lights red and blue => magenta
        pwmR.ChangeDutyCycle(100)
        pwmG.ChangeDutyCycle(0)
        pwmB.ChangeDutyCycle(100)
        time.sleep(t)

        # 绿灯，蓝灯全亮，红灯全暗（青色）
        # lights green and blue => aqua
        pwmR.ChangeDutyCycle(0)
        pwmG.ChangeDutyCycle(100)
        pwmB.ChangeDutyCycle(100)
        time.sleep(t)

        # 红灯，绿灯，蓝灯全亮（白色）
        # lights red, green and blue => white
        pwmR.ChangeDutyCycle(100)
        pwmG.ChangeDutyCycle(100)
        pwmB.ChangeDutyCycle(100)
        time.sleep(t)

except KeyboardInterrupt:
    pass

pwmR.stop()
pwmG.stop()
pwmB.stop()

gpio.cleanup()