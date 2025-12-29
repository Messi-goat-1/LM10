# PCAP Analyzer – Simple Explanation

This document explains the main functions used to analyze
a PCAP file in a simple and clear way.

---

## GetFileHandle(filePath)

This function is purely technical.

What it does:
- Opens a PCAP file stored on disk (offline file).
- Prepares the file for reading using the `pcap` library.
- Returns a file handle that can be used for packet processing.

NOTE:  
This function does not analyze data.
It only prepares the file for reading.

---

## RunFullAnalysis(handle)

This function is the **analysis engine**.

What it does:
- Reads network packets one by one.
- Extracts packet details such as:
  - Source and destination IP addresses.
  - Ports.
  - Payload data.

Current behavior:
- It is programmed to print only the first two packets.
- This is used as a test to verify the analysis pipeline.

TODO:  
Process all packets, not just the first two.

---

## AnalyzePCAP(fileID, filePath)

This is the main orchestration function.

What it does:
1. Calls `GetFileHandle` to open the PCAP file.
2. Passes the handle to `RunFullAnalysis`.
3. Controls the full analysis flow.

NOTE:  
This function is called directly by the server after the file
is fully assembled.

---

## Summary

- `GetFileHandle` opens the PCAP file.
- `RunFullAnalysis` analyzes packet data.
- `AnalyzePCAP` connects everything together.

This design separates file handling from analysis logic,
making the code easier to extend and maintain.
