import java.util.Scanner;

public class Sum {
    public static void main(String[] args) {
        final Scanner scanner = new Scanner(System.in);
        final int a = scanner.nextInt();
        final int b = scanner.nextInt();
        scanner.close();
        System.out.println(a + b);
    }
}