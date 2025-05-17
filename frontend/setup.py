from setuptools import setup, find_packages

setup(
    name="protopy",
    version="0.1",
    packages=find_packages(),
    install_requires=[
        'grpcio==1.71.0',
        'grpcio-tools==1.71.0',
    ],
)