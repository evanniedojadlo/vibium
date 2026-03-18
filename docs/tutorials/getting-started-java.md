# Getting Started with Vibium (Java)

A complete beginner's guide. No prior experience required.

---

## What You'll Build

A program that opens a browser, visits a website, takes a screenshot, and clicks a link. All in about 10 lines of code.

---

## Step 1: Install Java

Vibium requires Java 11 or higher (we recommend the latest LTS). Check if you have it:

```bash
java --version
```

If you see a version number (`11.0.0` or higher), skip to Step 2.

### macOS

```bash
# Install Homebrew (if you don't have it)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Java
brew install openjdk
```

### Windows

```powershell
winget install EclipseAdoptium.Temurin.21.JDK
```

### Linux

```bash
# Ubuntu/Debian
sudo apt-get install openjdk-21-jdk

# Or use your distro's package manager
```

---

## Step 2: Create a Project Folder

### Gradle (recommended)

```bash
mkdir my-first-bot
cd my-first-bot
```

Create a `build.gradle` file:

```groovy
plugins {
    id 'application'
}

repositories {
    mavenCentral()
}

dependencies {
    implementation 'com.vibium:vibium:26.3.17'
}

application {
    mainClass = 'Hello'
}
```

### Maven

```bash
mkdir my-first-bot
cd my-first-bot
```

Create a `pom.xml` file:

```xml
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>my-first-bot</artifactId>
    <version>1.0</version>

    <dependencies>
        <dependency>
            <groupId>com.vibium</groupId>
            <artifactId>vibium</artifactId>
            <version>26.3.17</version>
        </dependency>
    </dependencies>
</project>
```

### No build tool

Download the JAR directly and compile with `javac`:

```bash
mkdir my-first-bot
cd my-first-bot

# Download the vibium JAR (includes the vibium binary for your platform)
curl -LO https://repo1.maven.org/maven2/com/vibium/vibium/26.3.17/vibium-26.3.17.jar
```

---

## Step 3: Write Your First Program

Create a file called `Hello.java`:

```java
import com.vibium.Vibium;
import com.vibium.Browser;
import com.vibium.Page;
import com.vibium.Element;
import java.nio.file.Files;
import java.nio.file.Path;

public class Hello {
    public static void main(String[] args) throws Exception {
        // Launch a browser (you'll see it open!)
        Browser bro = Vibium.start();
        Page vibe = bro.page();

        // Go to a website
        vibe.go("https://example.com");
        System.out.println("Loaded example.com");

        // Take a screenshot
        byte[] png = vibe.screenshot();
        Files.write(Path.of("screenshot.png"), png);
        System.out.println("Saved screenshot.png");

        // Find and click the link
        Element link = vibe.find("a");
        System.out.println("Found link: " + link.text());
        link.click();
        System.out.println("Clicked!");

        // Close the browser
        bro.stop();
        System.out.println("Done!");
    }
}
```

---

## Step 4: Run It

### Gradle

```bash
# Place Hello.java in src/main/java/
mkdir -p src/main/java
mv Hello.java src/main/java/

gradle run
```

### Maven

```bash
# Place Hello.java in src/main/java/
mkdir -p src/main/java
mv Hello.java src/main/java/

mvn compile exec:java -Dexec.mainClass=Hello
```

### No build tool

```bash
javac -cp vibium-26.3.17.jar Hello.java
java -cp .:vibium-26.3.17.jar Hello
```

You should see:
1. A Chrome window open
2. example.com load
3. The browser click "More information..."
4. The browser close

Check your folder - there's now a `screenshot.png` file!

---

## What Just Happened?

| Line | What It Does |
|------|--------------|
| `Vibium.start()` | Opens Chrome, returns a Browser |
| `bro.page()` | Gets the default page (tab) |
| `vibe.go(url)` | Navigates to a URL |
| `vibe.screenshot()` | Captures the page as PNG bytes |
| `vibe.find(selector)` | Finds an element by CSS selector |
| `link.text()` | Gets the element's text content |
| `link.click()` | Clicks the element |
| `bro.stop()` | Closes the browser |

---

## Next Steps

**Hide the browser** (run headless):
```java
import com.vibium.types.StartOptions;

Browser bro = Vibium.start(new StartOptions().headless(true));
```

**Use JavaScript instead:**
See [Getting Started (JavaScript)](getting-started-js.md) for the JS version.

**Use Python instead:**
See [Getting Started (Python)](getting-started-python.md) for the Python version.

**Let AI control the browser:**
See [Getting Started with MCP](getting-started-mcp.md) for Claude Code and Gemini CLI setup.

---

## Troubleshooting

### "command not found: java"

Java isn't installed or isn't in your PATH. Reinstall from [adoptium.net](https://adoptium.net).

### "package com.vibium does not exist"

Make sure the vibium JAR is on your classpath. If using Gradle/Maven, run the build command to download dependencies first.

### Browser doesn't open

Try running with headless mode disabled (it's disabled by default, but just in case):
```java
Browser bro = Vibium.start(new StartOptions().headless(false));
```

### Permission denied (Linux)

You might need to install dependencies for Chrome:
```bash
sudo apt-get install -y libgbm1 libnss3 libatk-bridge2.0-0 libdrm2 libxkbcommon0 libxcomposite1 libxdamage1 libxfixes3 libxrandr2 libasound2
```

---

## You Did It!

You just automated a browser with Java. The same techniques work for:
- Web scraping
- Testing websites
- Automating repetitive tasks
- Building AI agents that can browse the web

Questions? [Open an issue](https://github.com/VibiumDev/vibium/issues).
