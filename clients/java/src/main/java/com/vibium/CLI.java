package com.vibium;

import com.vibium.internal.BinaryResolver;

/**
 * CLI proxy that forwards commands to the bundled vibium binary.
 *
 * <p>Allows Java users to run vibium CLI commands directly:
 * <pre>
 * java -jar vibium.jar install      # downloads Chrome for Testing
 * java -jar vibium.jar paths        # shows binary paths
 * </pre>
 */
public final class CLI {

    private CLI() {}

    public static void main(String[] args) {
        String binaryPath;
        try {
            binaryPath = BinaryResolver.resolve();
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            System.exit(1);
            return;
        }

        try {
            String[] cmd = new String[args.length + 1];
            cmd[0] = binaryPath;
            System.arraycopy(args, 0, cmd, 1, args.length);

            ProcessBuilder pb = new ProcessBuilder(cmd);
            pb.inheritIO();
            Process process = pb.start();
            int exitCode = process.waitFor();
            System.exit(exitCode);
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            System.exit(1);
        }
    }
}
