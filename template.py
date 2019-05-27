#!/usr/bin/env python

from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException

def chrome_driver():
    return webdriver.Chrome()

